//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	postID := "36387330-1cdb-4fc3-8ac4-ecaea768d7fb"
	err = client.Posts.Delete(ctx, postID)
	if err != nil {
		log.Fatal(err)
	}
}
