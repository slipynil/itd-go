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

	feed := client.Posts.NewFeed(types.FeedTabPopular, 10)

	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	for range timer.C {
		if !feed.HasMore() {
			break
		}

		posts, err := feed.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, post := range posts {
			pp.Println(post)
		}
	}
}
