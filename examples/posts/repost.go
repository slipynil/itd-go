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
	repostID := "3fb687f9-761f-44aa-847f-839ea560af5c"
	content := "а вы знали, что если вы сделаете репост и удалите, то потом не сможете сделать заново репост к тому посту"
	post, err := client.Posts.Repost(ctx, repostID, content)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(post)
}
