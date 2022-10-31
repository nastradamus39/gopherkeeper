package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gophkeeper/gopherkeeper/proto"
)

func registerSecrets(app *AppContext) {
	var secretCmd = &cobra.Command{
		Use:   "secret",
		Short: "Секреты",
	}

	// Список секретов
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Список секретов",
		Run: func(cmd *cobra.Command, args []string) {
			secrets, _ := app.grpcClient.SecretsListHandler(context.Background(), &proto.SecretsListRequest{})

			fmt.Println(secrets)
		},
	}

	// Добавление секрета
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Добавление нового секрета",
		Run: func(cmd *cobra.Command, args []string) {
			secret := proto.Secret{
				Login:   cmd.Flag("login").Value.String(),
				Comment: cmd.Flag("comment").Value.String(),
				Card:    cmd.Flag("card").Value.String(),
			}

			_, err := app.grpcClient.SecretsAddHandler(context.Background(), &proto.SecretsAddRequest{Secret: &secret})
			if err != nil {
				return
			}

			fmt.Println("Сохранено")
		},
	}
	addCmd.PersistentFlags().String("comment", "", "Комментарий к секрету")
	addCmd.PersistentFlags().String("login", "", "Логин")
	addCmd.PersistentFlags().String("card", "", "Данные банковской карты")

	secretCmd.AddCommand(listCmd)
	secretCmd.AddCommand(addCmd)

	app.rootCmd.AddCommand(secretCmd)
}
