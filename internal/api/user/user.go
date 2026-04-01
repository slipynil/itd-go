package user

import (
	"context"

	"github.com/slipynil/itd-go/internal/transport"
)

type User struct {
	transport *transport.Client
}

func New(t *transport.Client) *User {
	return &User{transport: t}
}

func (u *User) Get(ctx context.Context, id string) {

}
