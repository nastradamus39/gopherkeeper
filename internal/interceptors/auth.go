package interceptors

import (
	"context"
	"fmt"

	"gophkeeper/internal/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Необходима авторизация")
	}

	u, err := db.Repositories().Users.FindByToken(md["authorization"][0])

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Необходима авторизация")
	}

	ctx = context.WithValue(ctx, "user", u)

	resp, err = handler(ctx, req)

	// выполняем действия после вызова метода
	fmt.Println("выполняем действия после вызовом метода")

	return resp, err
}
