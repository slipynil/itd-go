package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
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

	feed := client.Posts.NewFeed(ctx, types.FeedTabPopular, 1)
	posts, err := feed.Next()
	if err != nil {
		log.Fatal(err)
	}

	firstPost := posts[0]

	pp.Printf("Автор поста: %s\n", firstPost.Author.DisplayName)
	pp.Printf("Контент: %s\n", firstPost)

	iterator := client.Comments.NewCommentList(ctx, firstPost.ID, 1)

	for iterator.HasMore() {
		comments, err := iterator.Next()
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
