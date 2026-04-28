//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
)

func main() {
	ctx := context.Background()
	cfg := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
	}

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

	content := "Опрос для разработчиков! Какой язык выбирают разработчики? 🚀"
	builder := types.NewPost(content).
		Bold("Опрос").
		Italic("разработчиков"). // Применится к обоим вхождениям слова
		Link("язык", "https://go.dev")

	// Создаём пост с опросом
	post, err := client.Posts.CreateWithPoll(ctx, builder, &poll)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Пост с опросом создан:")
	pp.Println(post)
}
