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

	postID := "c36ae616-765f-4119-8380-5fd8080df2d0"
	comment, err := client.Comments.CreateComment(ctx, postID, "привет", nil)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(comment)
}
