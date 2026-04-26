//go:build ignore

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

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	cfg := itdgo.Config{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
	}
	ctx := context.Background()

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
