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
	repostID := "3fb687f9-761f-44aa-847f-839ea560af5c"
	content := "а вы знали, что если вы сделаете репост и удалите, то потом не сможете сделать заново репост к тому посту"
	post, err := client.Posts.Repost(ctx, repostID, content)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(post)
}
