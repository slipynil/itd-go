//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	ctx := context.Background()
	cfg := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
		RetryDelay:   4 * time.Second,
	}

	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	iter, err := client.Search.NewHashtagPosts("nowkie", 5)
	if err != nil {
		log.Fatal(err)
	}

	for iter.HasMore() {
		posts, err := iter.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		pp.Println(posts)
	}
}
