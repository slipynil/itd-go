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
	userAgent := os.Getenv("USER_AGENT")
	token := os.Getenv("REFRESH_TOKEN")
	cfg := itdgo.Config{
		UserAgent:    userAgent,
		RefreshToken: token,
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	filePaths := []string{
		"/home/user/Pictures/cat.png",
		"/home/user/Pictures/cat.webp",
	}

	post, err := client.Posts.Create(ctx, "cat meme", filePaths...)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(post)
}
