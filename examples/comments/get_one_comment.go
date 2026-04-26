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

	commentIterator := client.Comments.NewCommentList(ctx, postID, 100)

	for commentIterator.HasMore() {
		comments, err := commentIterator.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, comment := range comments {
			if comment.ID == "dd000ea9-2268-4b5b-9be4-abeeef473702" {
				pp.Println(comment)
			}
		}

	}

}
