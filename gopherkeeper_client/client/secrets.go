package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/proto"

	"gophkeeper/internal/db"

	"github.com/spf13/cobra"
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
			fmt.Println("secrets list...")

			secrets, _ := app.grpcClient.SecretsListHandler(context.Background(), &proto.SecretsListRequest{})

			fmt.Println(secrets)
		},
	}

	// Добавление секрета
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Добавление нового секрета",
		Run: func(cmd *cobra.Command, args []string) {
			secret := db.Secret{
				Persist: false,
				Login:   cmd.Flag("login").Value.String(),
				Comment: cmd.Flag("comment").Value.String(),
				Card:    cmd.Flag("card").Value.String(),
			}

			fmt.Println(secret)
		},
	}
	addCmd.PersistentFlags().String("comment", "", "Комментарий к секрету")
	addCmd.PersistentFlags().String("login", "", "Логин")
	addCmd.PersistentFlags().String("card", "", "Данные банковской карты")

	// Обновление секрета
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Обновление существующего секрета",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("update secret...")
		},
	}
	updateCmd.PersistentFlags().String("comment", "", "Комментарий к секрету")
	updateCmd.PersistentFlags().String("login", "", "Логин")
	updateCmd.PersistentFlags().String("card", "", "Данные банковской карты")

	secretCmd.AddCommand(listCmd)
	secretCmd.AddCommand(addCmd)
	secretCmd.AddCommand(updateCmd)

	app.rootCmd.AddCommand(secretCmd)
}
