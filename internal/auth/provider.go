package auth

import (
	"context"
	"net/http"
)

type Provider interface {
	GetCookieJar() http.CookieJar
	GetAccessToken(ctx context.Context) (string, error)
	GetUserID() string
	IsAuthenticated() bool
}
