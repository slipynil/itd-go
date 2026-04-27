//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	ctx := context.Background()
	cfg := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	username := "nowkie"
	limit := 100
	pp.Printf("%v подписчиков пользователя %s:\n\n", limit, username)

	users, err := client.User.GetFollowers(ctx, "nowkie", limit)
	if err != nil {
		log.Fatalf("Ошибка получения подписчиков: %v", err)
	}
	for _, user := range users {
		pp.Printf("Пользователь: %s (%s)\n", user.Username, user.DisplayName)
	}
}
