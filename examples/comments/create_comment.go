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
