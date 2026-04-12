package root

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

const defaultTimeout = time.Second * 30

func CreateHttpClient(cfg Config) (*http.Client, error) {
	// Создаем cookie jar для refresh token
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании cookie jar: %w", err)
	}

	// Создаем cookie с refresh token
	cookie := &http.Cookie{
		Name:   "refresh_token",
		Value:  cfg.RefreshToken,
		Domain: cfg.Domain,
		Path:   "/api",
	}

	// Добавляем cookie в jar
	u, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("неверный формат baseURL: %w", err)
	}
	jar.SetCookies(u, []*http.Cookie{cookie})

	// Строим User-Agent
	userAgent := cfg.UserAgent

	// Устанавливаем timeout
	var timeout time.Duration
	if cfg.Timeout != 0 {
		timeout = cfg.Timeout
	} else {
		timeout = defaultTimeout
	}

	t, err := defaultAuthTransport()
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании базового transport: %w", err)
	}

	// Создаем HTTP клиент с cookie jar и базовыми заголовками
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: timeout,
		Transport: &headerTransport{
			base: t,
			headers: map[string]string{
				"Origin":     cfg.Url,
				"Referer":    cfg.Url + "/",
				"User-Agent": userAgent,
			},
		},
	}
	return httpClient, nil
}

// defaultAuthTransport создает базовый transport для auth операций
func defaultAuthTransport() (*http.Transport, error) {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       90 * time.Second,
	}
	if err := http2.ConfigureTransport(t); err != nil {
		return nil, fmt.Errorf("ошибка в поддержке HTTP/2.0: %w", err)
	}
	return t, nil
}

// headerTransport добавляет статические заголовки к запросам
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

// реализация интерфейса http.RoundTripper для http.Client
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.base.RoundTrip(req)
}
