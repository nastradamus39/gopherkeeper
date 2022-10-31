package main

import (
	_ "github.com/lib/pq"

	"bytes"
	"embed"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"

	"gophkeeper/gopherkeeper"
	"gophkeeper/gopherkeeper/proto"
	"gophkeeper/internal/db"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/interceptors"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	r := Router()

	// Logger
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer flog.Close()

	log.SetOutput(flog)

	// Переменные окружения в конфиг
	err = env.Parse(&gopherkeeper.Cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Инициализация подключения к бд
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	// миграции
	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
		return
	}
	if err = goose.Up(gopherkeeper.DB.DB, "migrations"); err != nil {
		log.Fatal(err)
		return
	}

	// получаем запрос gRPC
	go func() {
		log.Printf("Starting grpc")
		// определяем порт для grpc
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}

		// создаём gRPC-сервер без зарегистрированной службы
		creds, err := credentials.NewServerTLSFromFile(certificates())
		c := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptors.AuthInterceptor))

		// регистрируем сервис
		proto.RegisterSecretsServer(c, proto.UnimplementedSecretsServer{})
		if err := c.Serve(listen); err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Сервер gRPC начал работу")
	}()

	// запускаем сервер
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Printf("Не удалось запустить сервер. %s", err)
		return
	}
}

func Router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Post("/user/register", handlers.RegisterHTTPHandler)
	r.Post("/user/login", handlers.LoginHTTPHandler)

	// закрытые авторизацией эндпоинты
	r.Mount("/user", privateRouter())

	return r
}

// privateRouter Роутер для закрытых авторизацией эндпоинтов
func privateRouter() http.Handler {
	r := chi.NewRouter()

	return r
}

func certificates() (certPath string, keyPath string) {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certFile, err := os.OpenFile("./cert", os.O_CREATE|os.O_WRONLY, 0777)
	certFile.Write(certPEM.Bytes())
	certFile.Close()

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	keyFile, err := os.OpenFile("./key", os.O_CREATE|os.O_WRONLY, 0777)
	keyFile.Write(privateKeyPEM.Bytes())
	keyFile.Close()

	certPath = "cert"
	keyPath = "key"

	return
}
