//go:build ignore

package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
	"github.com/slipynil/itd-go/types"
)

func main() {
	go http.ListenAndServe(":6060", nil)
	ctx := context.Background()
	cfg := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
		RetryDelay:   5 * time.Second,
	}

	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	iter := client.Posts.NewFeed(types.FeedTabPopular, 20)
	var n int

	for iter.HasMore() {
		posts, err := iter.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		n += len(posts)
		pp.Println(n)
	}
	pp.Println("end successfull")
}
