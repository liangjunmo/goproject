package pagination

import (
	"github.com/spf13/cast"
)

type Request interface {
	Paginate(totalRecords interface{}) Pagination
}

type DefaultRequest struct {
	Page     uint32 `form:"page" json:"page"`
	Capacity uint32 `form:"capacity" json:"capacity"`
}

func (req DefaultRequest) Paginate(totalRecords interface{}) Pagination {
	return Paginate(req.Page, 1, req.Capacity, 10, cast.ToUint32(totalRecords))
}
