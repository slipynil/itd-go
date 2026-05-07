//go:build ignore

package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/k0kubun/pp"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	go http.ListenAndServe(":6060", nil)
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
		case err := <-errs:
			log.Fatal(err)
		case notification := <-stream:
			pp.Println(notification)
		}
	}
}
