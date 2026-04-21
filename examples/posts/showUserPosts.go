package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/joho/godotenv"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	// Загружаем переменные окружения
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

	// Создаём итератор для постов пользователя
	username := "nowkie"
	iterator := client.Posts.NewUserPosts(ctx, username, 20)

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
