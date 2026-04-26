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

	postID := "a7b2512f-7ee7-4d7f-8224-8d25e94bf0ed"

	// Создаём простой текстовый комментарий
	filePath := "/home/user/Pictures/cat.gif"
	comment, err := client.Comments.CreateComment(ctx, postID, "Fakemink Night, Blooming Jasmine", filePath)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Комментарий создан:")
	pp.Println(comment)
}
