package types

import "unicode"

// PostBuilder предоставляет fluent API для создания постов с форматированием текста.
// Поддерживает жирный, курсив, подчёркивание, зачёркивание, спойлер, моноширинный текст и ссылки.
//
// Форматирование применяется ко всем вхождениям целых слов в тексте.
// Слово определяется как последовательность букв и цифр, окружённая не-буквенными символами.
//
// Пример:
//
//	builder := types.NewPost("Go is awesome! I love Go programming.").
//		Bold("Go").        // Выделит оба вхождения "Go"
//		Italic("awesome")  // Выделит "awesome"
type PostBuilder struct {
	// Content - текстовое содержимое поста
	Content string
	// Spans - массив спанов (форматирования) в тексте
	Spans []Span
}

// NewPost создаёт новый билдер для поста с указанным содержимым.
func NewPost(content string) *PostBuilder {
	return &PostBuilder{Content: content}
}

// Bold применяет жирное форматирование ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Bold(text string) *PostBuilder {
	return b.addSpan("bold", text, "")
}

// Italic применяет курсивное форматирование ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Italic(text string) *PostBuilder {
	return b.addSpan("italic", text, "")
}

// Link создаёт ссылку для всех вхождений целого слова в тексте.
// Параметры:
//   - text: текст для поиска и преобразования в ссылку
//   - url: URL, на который будет вести ссылка
//
// Если слово не найдено или найдено как часть другого слова, ссылка не создаётся.
func (b *PostBuilder) Link(text string, url string) *PostBuilder {
	return b.addSpan("link", text, url)
}

// Underline применяет подчёркивание ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Underline(text string) *PostBuilder {
	return b.addSpan("underline", text, "")
}

// Strike применяет зачёркивание ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Strike(text string) *PostBuilder {
	return b.addSpan("strike", text, "")
}

// Spoiler применяет спойлер ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Spoiler(text string) *PostBuilder {
	return b.addSpan("spoiler", text, "")
}

// Monospace применяет моноширинное форматирование ко всем вхождениям целого слова в тексте.
// Если слово не найдено или найдено как часть другого слова, форматирование не применяется.
func (b *PostBuilder) Monospace(text string) *PostBuilder {
	return b.addSpan("monospace", text, "")
}

func (b *PostBuilder) addSpan(spanType, text, url string) *PostBuilder {
	contentRunes := []rune(b.Content)
	textRunes := []rune(text)
	textLen := len(textRunes)

	// Ищем все вхождения целых слов
	for i := range len(contentRunes) - textLen + 1 {
		// Проверяем совпадение текста в позиции i
		match := true
		for j := range textLen {
			if contentRunes[i+j] != textRunes[j] {
				match = false
				break
			}
		}

		if !match {
			continue
		}

		// Проверяем границы слова
		// До: либо начало строки, либо не буква/цифра
		beforeOk := i == 0 || !isWordChar(contentRunes[i-1])
		// После: либо конец строки, либо не буква/цифра
		afterOk := i+textLen == len(contentRunes) || !isWordChar(contentRunes[i+textLen])

		if beforeOk && afterOk {
			span := Span{
				Type:   spanType,
				Offset: i,
				Length: textLen,
			}
			if url != "" {
				span.URL = url
			}
			b.Spans = append(b.Spans, span)
		}
	}

	return b
}

// isWordChar проверяет, является ли символ частью слова (буква или цифра)
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
