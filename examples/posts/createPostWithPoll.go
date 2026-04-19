package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	cfg := itdgo.Config{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
	}
	ctx := context.Background()

	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём опрос с вопросом и вариантами ответов
	poll := types.PollRequest{
		Question: "Какой язык программирования вам нравится больше?",
		Options: []types.PollOptionRequest{
			{Text: "Go"},
			{Text: "Python"},
			{Text: "JavaScript"},
			{Text: "Rust"},
		},
		MultipleChoice: false, // Можно выбрать только один вариант
	}

	content := "Опрос для разработчиков! 🚀"

	// Создаём пост с опросом
	post, err := client.Posts.CreateWithPoll(ctx, content, &poll)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Пост с опросом создан:")
	pp.Println(post)
}
