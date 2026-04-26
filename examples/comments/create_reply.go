//go:build ignore

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
	godotenv.Load()
	token := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	cfg := itdgo.Config{
		RefreshToken: token,
		UserAgent:    userAgent,
	}

	ctx := context.Background()

	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	postID := "c36ae616-765f-4119-8380-5fd8080df2d0"

	// Получаем первую страницу комментариев
	iter := client.Comments.NewCommentList(ctx, postID, 10)
	if !iter.HasMore() {
		log.Fatal("Нет комментариев для ответа")
	}

	comments, err := iter.Next(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(comments) == 0 {
		log.Fatal("Нет комментариев для ответа")
	}

	// Отвечаем на первый комментарий
	firstComment := comments[0]

	// Создаём простой текстовый ответ
	reply, err := client.Comments.CreateReply(
		ctx,
		firstComment.ID,
		firstComment.Author.ID,
		"Согласен! 💯",
	)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Ответ создан:")
	pp.Println(reply)

	// Для создания ответа с файлами используйте:
	// reply, err := client.Comments.CreateReply(
	//     ctx,
	//     firstComment.ID,
	//     firstComment.Author.ID,
	//     "Вот скриншот",
	//     "/path/to/screenshot.png",
	// )
	// if err != nil {
	//     log.Fatal(err)
	// }
	// pp.Println(reply)
}
