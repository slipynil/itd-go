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
		UserAgent:    os.Getenv("USERAGENT"),
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	iter := client.Notifications.NewIterator(10)

	for iter.HasMore() {
		notifications, err := iter.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		for _, notification := range notifications {
			pp.Println(notification)
		}
	}
}
