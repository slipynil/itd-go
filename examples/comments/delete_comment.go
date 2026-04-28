//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
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

	commentID := "your-comment-id-here"
	err = client.Comments.Delete(ctx, commentID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Комментарий успешно удалён")
}
