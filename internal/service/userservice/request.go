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

type GetUserByUIDRequest struct {
	UID uint32
}

type GetUserByUsernameRequest struct {
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
