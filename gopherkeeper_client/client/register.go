package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"

	"gophkeeper/internal/proto"

	"github.com/spf13/cobra"
)

func registerRegistration(app *AppContext) {
	cmd := &cobra.Command{
		Use:   "register <login> <password>",
		Short: "Регистрирует пользователя на сервере",
		Run: func(cmd *cobra.Command, args []string) {
			login, pass := parseLoginPassword(cmd, args)

			_, err := app.grpcClient.RegisterHandler(context.Background(), &proto.RegisterRequest{
				Login:    login,
				Password: pass,
			})

			// Разбираем ошибку, если такая вернулась
			if err != nil {
				if e, ok := status.FromError(err); ok {
					fmt.Printf("%s\n", e.Message())
				} else {
					fmt.Println("Неизвестная ошибка")
				}
				return
			}

			fmt.Println("Пользователь зарегистрирован")
		},
	}

	app.rootCmd.AddCommand(cmd)
}
