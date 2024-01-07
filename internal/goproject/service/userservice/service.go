package userservice

import (
	"context"

	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type Service interface {
	List(ctx context.Context, cmd ListCommand) (pagination.Pagination, []User, error)
	Search(ctx context.Context, cmd SearchCommand) (map[uint32]User, error)
	Get(ctx context.Context, cmd GetCommand) (User, error)
	GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (User, error)
	Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

type ListCommand struct {
	pagination.Request
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
