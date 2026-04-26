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

	ctx := context.Background()

	cfg := itdgo.Config{
		RefreshToken: token,
		UserAgent:    userAgent,
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("USER INFO")
	user, err := client.User.Me(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(user)

}
