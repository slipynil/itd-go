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
	postID := "ccd1e015-0b3c-4fec-a48b-744b269742e6"
	pollOptionID := "4e53e76f-1dbf-4cd0-9deb-bb007b26add1"
	polls, err := client.Posts.Vote(ctx, postID, pollOptionID)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(polls)
}
