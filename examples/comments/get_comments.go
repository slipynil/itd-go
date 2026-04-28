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

	postID := "a7b2512f-7ee7-4d7f-8224-8d25e94bf0ed"
	iter := client.Comments.NewCommentList(postID, 10)

	if iter.HasMore() {
		comments, err := iter.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, comment := range comments {
			pp.Println(comment)
		}
	}
}
