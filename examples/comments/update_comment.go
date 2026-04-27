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

	commentID := "dd000ea9-2268-4b5b-9be4-abeeef473702"
	newContent := "#привет #dota"

	result, err := client.Comments.Update(ctx, commentID, newContent)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(result)
}
