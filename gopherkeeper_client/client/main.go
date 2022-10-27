package main

import (
	_ "github.com/lib/pq"

	"crypto/tls"
	"fmt"
	"log"

	"gophkeeper/gopherkeeper_client/client/interceptors"
	"gophkeeper/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spf13/cobra"
)

const TokenPath = ".gph_key"

type AppContext struct {
	rootCmd        *cobra.Command
	grpcConnection *grpc.ClientConn
	grpcClient     proto.SecretsClient
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gphc",
		Short: "gopher keeper client",
	}

	appCtx := initAppContext(rootCmd)

	defer func(grpcConnection *grpc.ClientConn) {
		err := grpcConnection.Close()
		if err != nil {

		}
	}(appCtx.grpcConnection)

	registerRegistration(appCtx) // команда регистрации
	registerAuth(appCtx)         // команда авторизации
	registerSecrets(appCtx)      // список секретов

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func initAppContext(rootCmd *cobra.Command) *AppContext {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(
		":3200",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		grpc.WithUnaryInterceptor(interceptors.AuthInterceptor),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := proto.NewSecretsClient(conn)

	return &AppContext{
		rootCmd:        rootCmd,
		grpcConnection: conn,
		grpcClient:     client,
	}
}
