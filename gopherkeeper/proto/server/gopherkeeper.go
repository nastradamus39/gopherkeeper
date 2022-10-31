package server

import (
	"context"
	"fmt"

	"gophkeeper/gopherkeeper/proto"
	"gophkeeper/internal/db"
	"gophkeeper/internal/handlers"
)

// GopherKeeper поддерживает все необходимые методы сервера.
type GopherKeeper struct {
	// Нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedSecretsServer
}

// RegisterHandler Проверяет логин/пароль и возвращает токен
func (s *GopherKeeper) RegisterHandler(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := db.User{
		Login:    in.Login,
		Password: in.Password,
	}

	err := handlers.RegisterHandler(&user)

	if err != nil {
		return &proto.RegisterResponse{}, err
	}

	return &proto.RegisterResponse{Token: user.Token}, nil
}

// AuthorizeHandler Проверяет логин/пароль и возвращает токен
func (s *GopherKeeper) AuthorizeHandler(ctx context.Context, in *proto.AuthorizeRequest) (*proto.AuthorizeResponse, error) {
	token, err := handlers.LoginGRPCHandler(in.Login, in.Password)

	if err != nil {
		return &proto.AuthorizeResponse{}, err
	}

	return &proto.AuthorizeResponse{
		Token: token,
	}, nil
}

// SecretsListHandler Список секретов пользователя
func (s *GopherKeeper) SecretsListHandler(ctx context.Context, in *proto.SecretsListRequest) (*proto.SecretsListResponse, error) {
	fmt.Println("SecretsListHandler")

	var secrets []*proto.Secret

	secrets = append(secrets, &proto.Secret{
		Login:   "test",
		Comment: "sdasda",
		Card:    "cvjhzcovh",
	})

	return &proto.SecretsListResponse{Secrets: secrets}, nil
}

// SecretsAddHandler Добавляет секрет пользователя
func (s *GopherKeeper) SecretsAddHandler(ctx context.Context, in *proto.SecretsAddRequest) (*proto.SecretsAddResponse, error) {
	fmt.Println("SecretsAddHandler")

	secret := db.Secret{
		Persist: false,
		Login:   in.Secret.Login,
		Comment: in.Secret.Comment,
		Card:    in.Secret.Card,
	}

	err := db.Repositories().Secrets.Save(secret)
	if err != nil {
		return nil, err
	}

	return &proto.SecretsAddResponse{}, nil
}
