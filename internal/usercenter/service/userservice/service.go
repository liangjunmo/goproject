package userservice

import (
	"context"
)

type Service interface {
	Search(ctx context.Context, cmd SearchCommand) (map[uint32]User, error)
	Get(ctx context.Context, cmd GetCommand) (User, error)
	GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (User, error)
	Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

type SearchCommand struct {
	Uids     []uint32
	Username string
}

type GetCommand struct {
	UID uint32
}

type GetByUsernameCommand struct {
	Username string
}

type CreateCommand struct {
	Username string
	Password string
}

type ValidatePasswordCommand struct {
	Username string
	Password string
}
