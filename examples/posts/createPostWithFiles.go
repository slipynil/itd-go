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

	content := "Fakemink Night, Blooming Jasmine"

	// Указываем пути к файлам для загрузки
	// Файлы будут автоматически загружены и прикреплены к посту
	// filePaths := []string{
	// 	"/home/user/Downloads/firefox/320kg.jpg",
	// 	"/home/user/Pictures/tuntun.jpg",
	// }

	filePath := "/home/user/Downloads/music/fakemink_night_blooming_Jasmine.mp3"

	// Создаём пост с файлами
	post, err := client.Posts.Create(ctx, content, filePath)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println("Пост с файлами создан:")
	pp.Println(post)
	pp.Println("\nПрикреплённые файлы:")
	for i, attachment := range post.Attachments {
		pp.Printf("%d. %s (ID: %s)\n", i+1, attachment.URL, attachment.ID)
	}
}
