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

	ctx := context.Background()

	cfg := itdgo.Config{
		RefreshToken: token,
		UserAgent:    userAgent,
	}
	client, err := itdgo.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Обновляем профиль с автоматической загрузкой баннера
	profile := types.UpdateProfile{
		Bio:        "Новая биография",
		BannerPath: "/home/user/Pictures/fakemink-ig.png",
	}

	user, err := client.User.UpdateProfile(ctx, profile)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Профиль обновлён:")
	pp.Println(user)
}
