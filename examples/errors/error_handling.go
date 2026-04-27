//go:build ignore

package main

import (
	"context"
	"errors"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	itdgo "github.com/slipynil/itd-go"
	itderrors "github.com/slipynil/itd-go/errors"
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

	log.Println("=== Примеры обработки ошибок в itd-go SDK ===\n")

	// 1. Валидационные ошибки - пустой ID
	demonstrateEmptyIDError(ctx, client)

	// 2. Валидационные ошибки - пустой контент
	demonstrateEmptyContentError(ctx, client)

	// 3. Валидационные ошибки - опрос
	demonstratePollValidationErrors(ctx, client)

	// 4. API ошибки - 404 Not Found
	demonstrateNotFoundError(ctx, client)

	// 5. Валидация файлов
	demonstrateFileValidationErrors(ctx, client)

	// 6. Получение деталей API ошибки
	demonstrateAPIErrorDetails(ctx, client)
}

// 1. Демонстрация ошибки пустого ID
func demonstrateEmptyIDError(ctx context.Context, client *itdgo.Client) {
	log.Println("1. Проверка ошибки пустого ID поста:")

	_, err := client.Posts.Get(ctx, "")
	if err != nil {
		// Проверка через errors.Is() - рекомендуемый способ
		if errors.Is(err, itderrors.ErrEmptyPostID) {
			log.Println("   ✓ Обнаружена ошибка: ID поста не может быть пустым")
		} else {
			log.Printf("   ✗ Неожиданная ошибка: %v\n", err)
		}
	}
	log.Println()
}

// 2. Демонстрация ошибки пустого контента
func demonstrateEmptyContentError(ctx context.Context, client *itdgo.Client) {
	log.Println("2. Проверка ошибки пустого контента:")

	_, err := client.Posts.Create(ctx, "")
	if err != nil {
		if errors.Is(err, itderrors.ErrEmptyContent) {
			log.Println("   ✓ Обнаружена ошибка: контент или файлы обязательны")
		} else {
			log.Printf("   ✗ Неожиданная ошибка: %v\n", err)
		}
	}
	log.Println()
}

// 3. Демонстрация ошибок валидации опроса
func demonstratePollValidationErrors(ctx context.Context, client *itdgo.Client) {
	log.Println("3. Проверка ошибок валидации опроса:")

	// 3.1. Nil poll
	_, err := client.Posts.CreateWithPoll(ctx, "Опрос", nil)
	if err != nil {
		if errors.Is(err, itderrors.ErrNilPoll) {
			log.Println("   ✓ Обнаружена ошибка: poll не может быть nil")
		}
	}

	// 3.2. Недостаточно вариантов
	poll := &types.PollRequest{
		Question: "Вопрос?",
		Options: []types.PollOptionRequest{
			{Text: "Вариант 1"},
		},
	}
	_, err = client.Posts.CreateWithPoll(ctx, "Опрос", poll)
	if err != nil {
		if errors.Is(err, itderrors.ErrInsufficientPollOptions) {
			log.Println("   ✓ Обнаружена ошибка: в опросе должно быть минимум 2 варианта")
		}
	}
	log.Println()
}

// 4. Демонстрация API ошибки 404
func demonstrateNotFoundError(ctx context.Context, client *itdgo.Client) {
	log.Println("4. Проверка API ошибки 404 Not Found:")

	_, err := client.Posts.Get(ctx, "non-existent-post-id-12345")
	if err != nil {
		// Проверка через errors.Is() - работает благодаря Unwrap()
		if errors.Is(err, itderrors.ErrNotFound) {
			log.Println("   ✓ Обнаружена ошибка: пост не найден (404)")
		} else if errors.Is(err, itderrors.ErrUnauthorized) {
			log.Println("   ✓ Обнаружена ошибка: токен истёк (401)")
		} else {
			log.Printf("   ℹ Другая ошибка: %v\n", err)
		}
	}
	log.Println()
}

// 5. Демонстрация ошибок валидации файлов
func demonstrateFileValidationErrors(ctx context.Context, client *itdgo.Client) {
	log.Println("5. Проверка ошибок валидации файлов:")

	// 5.1. Неподдерживаемое расширение
	_, err := client.Posts.Create(ctx, "Пост с файлом", "/path/to/file.mp3")
	if err != nil {
		if errors.Is(err, itderrors.InvalidFileExtension) {
			log.Println("   ✓ Обнаружена ошибка: неподдерживаемое расширение файла")
		}
	}

	// 5.2. Файл без расширения
	_, err = client.Posts.Create(ctx, "Пост с файлом", "/path/to/file")
	if err != nil {
		if errors.Is(err, itderrors.NoFileExtension) {
			log.Println("   ✓ Обнаружена ошибка: файл без расширения")
		}
	}

	// 5.3. Слишком много файлов
	tooManyFiles := make([]string, 11)
	for i := range tooManyFiles {
		tooManyFiles[i] = "/path/to/file.png"
	}
	_, err = client.Posts.Create(ctx, "Пост с файлами", tooManyFiles...)
	if err != nil {
		if errors.Is(err, itderrors.TooManyFiles) {
			log.Println("   ✓ Обнаружена ошибка: превышен лимит файлов (макс. 10)")
		}
	}
	log.Println()
}

// 6. Демонстрация получения деталей API ошибки
func demonstrateAPIErrorDetails(ctx context.Context, client *itdgo.Client) {
	log.Println("6. Получение деталей API ошибки:")

	_, err := client.Posts.Get(ctx, "non-existent-post-id-12345")
	if err != nil {
		// Получение деталей через errors.As()
		var apiErr *itderrors.APIError
		if errors.As(err, &apiErr) {
			log.Printf("   ✓ Детали API ошибки:\n")
			log.Printf("     - Код: %s\n", apiErr.Code)
			log.Printf("     - Сообщение: %s\n", apiErr.Message)
			log.Printf("     - HTTP статус: %d\n", apiErr.StatusCode)

			// Проверка обёрнутой sentinel ошибки
			if apiErr.Err != nil {
				log.Printf("     - Тип ошибки: %v\n", apiErr.Err)
			}
		}
	}
	log.Println()
}
