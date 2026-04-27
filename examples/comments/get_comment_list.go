//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
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

	feed := client.Posts.NewFeed(types.FeedTabPopular, 1)
	posts, err := feed.Next(ctx)
	if err != nil {
		log.Fatal(err)
	}

	firstPost := posts[0]

	pp.Printf("Автор поста: %s\n", firstPost.Author.DisplayName)
	pp.Printf("Контент: %s\n", firstPost)

	iterator := client.Comments.NewCommentList(firstPost.ID, 1)

	for iterator.HasMore() {
		comments, err := iterator.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if len(comments) == 0 {
			pp.Println("--- НЕТ КОММЕНТАРИЕВ ---")
		}
		for _, comment := range comments {
			pp.Println("--- КОММЕНТАРИЙ ---")
			pp.Println(comment)
			pp.Println()
		}
	}
}
