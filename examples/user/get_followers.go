package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
)

func main() {

	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	token := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	ctx := context.Background()

	// Создаём конфигурацию и клиент ITD
	cfg := itdgo.Config{
		RefreshToken: token,
		UserAgent:    userAgent,
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
