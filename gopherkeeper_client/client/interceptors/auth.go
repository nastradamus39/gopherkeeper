package interceptors

import (
	"context"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tokenAuth struct {
	token string
}

// GetRequestMetadata Добавляет токен в заголовки
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func AuthInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	token, err := ioutil.ReadFile(".gph_key")
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Необходимо авторизоваться")
	}

	opts = append(opts, grpc.PerRPCCredsCallOption{Creds: tokenAuth{token: string(token)}})

	// вызываем RPC-метод
	err = invoker(ctx, method, req, reply, cc, opts...)

	return err
}
