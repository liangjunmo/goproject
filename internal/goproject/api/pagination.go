package api

import (
	"github.com/spf13/cast"

	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type Pagination struct {
	Page            uint32 `json:"page"`
	CapacityPerPage uint32 `json:"capacity_per_page"`
	TotalPages      uint32 `json:"total_pages"`
	TotalRecords    uint32 `json:"total_records"`
	Offset          int    `json:"-"`
	Limit           int    `json:"-"`
}

func (p Pagination) GetPage() uint32 {
	return p.Page
}

func (p Pagination) GetCapacityPerPage() uint32 {
	return p.CapacityPerPage
}

func (p Pagination) GetTotalPages() uint32 {
	return p.TotalPages
}

func (p Pagination) GetTotalRecords() uint32 {
	return p.TotalRecords
}

func (p Pagination) GetOffset() int {
	return p.Offset
}

func (p Pagination) GetLimit() int {
	return p.Limit
}

type PaginationRequest struct {
	Page     uint32 `form:"page" json:"page"`
	Capacity uint32 `form:"-" json:"-"`
}

func (req PaginationRequest) Paginate(totalRecords interface{}) pagination.Pagination {
	return pagination.Paginate(req.Page, 1, req.Capacity, 10, cast.ToUint32(totalRecords))
}
