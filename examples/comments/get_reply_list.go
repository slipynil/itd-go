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

	commentID := "b262e8c4-24ca-4f4f-82b0-dbcdacfecdb3"

	replies, err := client.Comments.ListReplies(ctx, commentID, 100)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(replies)
}
