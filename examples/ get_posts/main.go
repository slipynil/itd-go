package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	itd "github.com/slipynil/itd-go"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	cfg := itd.Config{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
	}
	ctx := context.Background()

	client, err := itd.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	feed := client.Posts.NewFeed(ctx, 1, "popular")
	posts, _ := feed.Next()
	pp.Println(posts)
}
