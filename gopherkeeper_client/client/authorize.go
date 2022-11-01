package main

import (
	"context"
	"fmt"
	"os"

	"gophkeeper/gopherkeeper/proto"
	"gophkeeper/internal/db"

	"github.com/spf13/cobra"

	"google.golang.org/grpc/status"
)

func registerAuth(app *AppContext) {
	cmd := &cobra.Command{
		Use:   "auth <login> <password>",
		Short: "Авторизует пользователя на сервере",
		Run: func(cmd *cobra.Command, args []string) {
			login, pass := parseLoginPassword(cmd, args)
			user := db.User{
				Login:    login,
				Password: pass,
			}

			response, err := app.grpcClient.AuthorizeHandler(context.Background(), &proto.AuthorizeRequest{
				Login:    user.Login,
				Password: user.Password,
			})

			// Разбираем ошибку, если такая вернулась
			if err != nil {
				if e, ok := status.FromError(err); ok {
					fmt.Printf("(%v) %s\n", e.Code(), e.Message())
				} else {
					fmt.Println("Неизвестная ошибка")
				}
			}

			// Аунтефицируем пользователя
			file, err := os.OpenFile(TokenPath, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println(err)
			}

			// Сохраняем токен
			_, err = file.Write([]byte(response.Token))
			if err != nil {
				fmt.Println(err)
			}
			err = file.Close()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Успешно")
		},
	}

	app.rootCmd.AddCommand(cmd)
}

func parseLoginPassword(cmd *cobra.Command, args []string) (login string, pass string) {
	if len(args) < 1 {
		fmt.Println("Передано слишком мало аргументов")
		_ = cmd.Help() // Help always returns no error
		os.Exit(1)
	}

	login = args[0]
	if login == "" {
		fmt.Println("Логин не может быть пустым")
		_ = cmd.Help() // Help always returns no error
		os.Exit(1)
	}

	pass = args[1]
	if pass == "" {
		fmt.Println("Пароль не может быть пустым")
		_ = cmd.Help() // Help always returns no error
		os.Exit(1)
	}

	return
}
