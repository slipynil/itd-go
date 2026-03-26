package client

import (
	"net"
	"net/http"
	"time"
)

// ITDurl - базовый URL для всех запросов к ITD API
const ITDurl = "https://xn--d1ah4a.com/api"

type client struct {
	httpClient *http.Client
	baseURL    string
}

// Конструктор клиента
// можно задать свой http.Client или
// использовать nil поумолчанию
func New(httpClient *http.Client) *client {

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: defaultTransport(), // дефолтный транспорт с настройками по умолчанию
			Timeout:   60 * time.Second,   // общий дедлайн на весь запрос
		}
	}

	return &client{
		httpClient: httpClient,
		baseURL:    ITDurl,
	}
}

// http.Transport - движок под капотом http.Client
// Отвечает за установку TCP-соединений, TLS-рукопожатия,
// connection pool, работу с proxy и выбор протокола
func defaultTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // максимум времени на установку TCP-соединения
			KeepAlive: 30 * time.Second, // включает TCP keepalive-пакеты каждые 30 секунд
		}).DialContext,
		ForceAttemptHTTP2:     true,             // поддержка HTTP/2
		TLSHandshakeTimeout:   10 * time.Second, // лимит на время TLS-хендшейка после установки TCP
		ResponseHeaderTimeout: 30 * time.Second, // время ожидания заголовков ответа после того, как полностью отправлен
		ExpectContinueTimeout: time.Second,

		// Настройка idle соединений
		MaxIdleConns:        100,              // максимум idle соединений
		MaxIdleConnsPerHost: 10,               // максимум idle-соейдинений на 1 хост
		MaxConnsPerHost:     0,                // общий лимит соединений (активных и idle) на хост
		IdleConnTimeout:     90 * time.Second, // время жизни idle соединений
	}
}
