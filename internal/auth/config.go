package auth

import "net/http"

type Config struct {
	Url        string
	HttpClient *http.Client
}
