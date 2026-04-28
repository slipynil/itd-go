//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	_ "github.com/joho/godotenv/autoload"
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

	// Создаём итератор для постов пользователя
	username := "nowkie"
	iterator := client.Posts.NewUserPosts(username, 20)

	fmt.Printf("Загрузка постов пользователя @%s...\n\n", username)

	pageNum := 1
	totalPosts := 0

	// Итерируемся по всем страницам
	for iterator.HasMore() {
		posts, err := iterator.Next(ctx)
		if err != nil {
			log.Fatalf("Ошибка получения постов: %v", err)
		}

		fmt.Printf("Страница %d: получено %d постов\n", pageNum, len(posts))

		for _, post := range posts {
			fmt.Printf("  - [%s] %s (лайков: %d)\n",
				post.CreatedAt.String(),
				truncate(post.Content, 50),
				post.LikesCount,
			)
		}

		totalPosts += len(posts)
		pageNum++
	}

	fmt.Printf("\nВсего загружено: %d постов\n", totalPosts)
}

func truncate(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen]) + "..."
}
