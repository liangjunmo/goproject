package accountservice

import (
	"context"
)

type Service interface {
	Login(ctx context.Context, cmd LoginCommand) (ticket string, failedCount uint32, err error)
	CreateToken(ctx context.Context, cmd CreateTokenCommand) (token string, err error)
	Authorize(ctx context.Context, cmd AuthorizeCommand) (*UserJwtClaims, error)
}

type LoginCommand struct {
	Username string
	Password string
}

type CreateTokenCommand struct {
	Ticket string
}

type AuthorizeCommand struct {
	Token string
}
