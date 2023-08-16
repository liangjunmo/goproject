package userservice

import (
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type ListUserCommand struct {
	PaginationRequest pagination.Request
}

type SearchUserCommand struct {
	Uids      []uint32
	Usernames []string
}

type GetUserCommand struct {
	Uid      uint32
	Username string
}

type CreateUserCommand struct {
	Username string
}
