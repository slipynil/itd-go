package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	godotenv.Load()
	token := os.Getenv("REFRESH_TOKEN")
	userAgent := os.Getenv("USER_AGENT")

	ctx := context.Background()

	cfg := itdgo.Config{
		RefreshToken: token,
		UserAgent:    userAgent,
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	username := "nowkie"
	if err := client.User.Follow(ctx, username); err != nil {
		log.Fatal(err)
	}
}
