//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	itdgo "github.com/slipynil/itd-go"
)

func main() {
	ctx := context.Background()

	log.Println("=== Демонстрация автоматической обработки rate limiting (429) ===\n")

	// Пример 1: Конфигурация по умолчанию (3 попытки, 1 секунда)
	log.Println("1. Конфигурация по умолчанию:")
	cfg1 := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
		// MaxRetries и RetryDelay не указаны - используются значения по умолчанию
	}

	client1, err := itdgo.New(ctx, cfg1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   - MaxRetries: 3 (по умолчанию)")
	log.Println("   - RetryDelay: 1s (по умолчанию)")
	log.Println("   - Exponential backoff: 1s → 2s → 4s\n")

	// Пример запроса, который может вызвать 429
	iter, err := client1.Search.NewHashtagPosts("nowkie", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   Выполняем запрос к API...")
	if iter.HasMore() {
		posts, err := iter.Next(ctx)
		if err != nil {
			log.Printf("   ✗ Ошибка: %v\n", err)
		} else {
			log.Printf("   ✓ Успешно получено %d постов\n", len(posts))
		}
	}
	log.Println()

	// Пример 2: Кастомная конфигурация
	log.Println("2. Кастомная конфигурация:")
	cfg2 := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
		MaxRetries:   5,                   // 5 попыток
		RetryDelay:   2 * time.Second,     // начальная задержка 2 секунды
	}

	client2, err := itdgo.New(ctx, cfg2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   - MaxRetries: 5")
	log.Println("   - RetryDelay: 2s")
	log.Println("   - Exponential backoff: 2s → 4s → 8s → 16s → 32s\n")

	// Пример 3: Отключение retry
	log.Println("3. Отключение retry логики:")
	cfg3 := itdgo.Config{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		UserAgent:    os.Getenv("USER_AGENT"),
		MaxRetries:   0, // отключить retry
	}

	client3, err := itdgo.New(ctx, cfg3)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("   - MaxRetries: 0 (retry отключен)")
	log.Println("   - При ошибке 429 запрос завершится немедленно без повторов\n")

	// Демонстрация поведения
	log.Println("=== Поведение при ошибке 429 ===")
	log.Println("Если API вернёт 429, вы увидите сообщения:")
	log.Println("  [itd-go] Rate limit exceeded (429), retry 1/3 after 1s")
	log.Println("  [itd-go] Rate limit exceeded (429), retry 2/3 after 2s")
	log.Println("  [itd-go] Rate limit exceeded (429), retry 3/3 after 4s")
	log.Println("\nSDK автоматически повторит запрос с увеличивающейся задержкой.")
	log.Println("Другие ошибки (401, 404, 5xx) возвращаются немедленно без повторов.")

	_ = client2
	_ = client3
}
