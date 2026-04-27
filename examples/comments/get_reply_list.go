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

	commentID := "b262e8c4-24ca-4f4f-82b0-dbcdacfecdb3"

	replies, err := client.Comments.ListReplies(ctx, commentID, 100)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(replies)
}
