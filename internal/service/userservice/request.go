package userservice

import (
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type ListUserRequest struct {
	PaginationRequest pagination.Request
}

type SearchUserRequest struct {
	Uids []uint32
}

type GetUserByUIDRequest struct {
	UID uint32
}

type CreateUserRequest struct {
	UID uint32
}
