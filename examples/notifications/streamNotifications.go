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
	stream, errs := client.Notifications.Stream(ctx)
	for {
		select {
		case n, ok := <-stream:
			if !ok {
				return
			}
			pp.Println(n)
		case err := <-errs:
			log.Println(err)
			return
		}
	}
}
