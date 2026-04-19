package transport

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/go-json-experiment/json"
	"github.com/slipynil/itd-go/types"
)

// Upload загружает файл на сервер ITD и возвращает информацию о вложении.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - path: путь к файлу на локальной файловой системе
//
// Возвращает Attachment с ID и URL загруженного файла или ошибку при проблемах с загрузкой.
func (c *Client) Upload(ctx context.Context, path string) (*types.Attachment, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := c.NewRequestMultipart(ctx, "POST", "/api/files/upload", body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.Attachment
	return &result, json.UnmarshalRead(resp.Body, &result)
}

// UploadFiles загружает несколько файлов и возвращает их ID.
// Параметры:
//   - ctx: контекст для управления временем жизни запроса
//   - filePaths: пути к файлам для загрузки
//
// Возвращает срез ID загруженных файлов или ошибку при проблемах с загрузкой.
func (c *Client) UploadFiles(ctx context.Context, filePaths ...string) ([]string, error) {
	if len(filePaths) == 0 {
		return nil, nil
	}

	attachmentIDs := make([]string, 0, len(filePaths))
	for _, path := range filePaths {
		attachment, err := c.Upload(ctx, path)
		if err != nil {
			return nil, err
		}
		attachmentIDs = append(attachmentIDs, attachment.ID)
	}
	return attachmentIDs, nil
}
