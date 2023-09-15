package userservice

import (
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type ListUserRequest struct {
	PaginationRequest pagination.Request
}

type SearchUserRequest struct {
	Uids      []uint32
	Usernames []string
}

type GetUserRequest struct {
	Uid      uint32
	Username string
}

type CreateUserRequest struct {
	Username string
	Password string
}

type ValidatePasswordRequest struct {
	Username string
	Password string
}
