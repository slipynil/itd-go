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

	filePaths := []string{
		"/home/user/Pictures/cat.png",
		"/home/user/Pictures/cat.webp",
	}

	// Демонстрация форматирования: Bold применится ко всем вхождениям слова "Go"
	builder := types.NewPost("Go is awesome! I love Go programming. Go Go Go!").
		Bold("Go").
		Italic("awesome").
		Underline("love")

	post, err := client.Posts.Create(ctx, builder, filePaths...)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(post)
}
