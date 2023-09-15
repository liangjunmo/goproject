package pagination

import "github.com/spf13/cast"

const (
	DefaultPageIndex uint32 = 1
	DefaultPageSize  uint32 = 10
)

type Pagination struct {
	PageIndex   uint32 `json:"page_index"`
	PageSize    uint32 `json:"page_size"`
	PageTotal   uint32 `json:"page_total"`
	ResultTotal uint32 `json:"result_total"`
	Offset      int    `json:"-"`
	Limit       int    `json:"-"`
}

type Request struct {
	PageIndex uint32 `form:"page_index" json:"page_index"`
	PageSize  uint32 `form:"page_index" json:"page_size"`
}

func (req Request) Paginate(_total interface{}) Pagination {
	if req.PageIndex < 1 {
		req.PageIndex = DefaultPageIndex
	}

	if req.PageSize < 1 {
		req.PageSize = DefaultPageSize
	}

	var (
		total     = cast.ToUint32(_total)
		pageTotal = total / req.PageSize
		offset    uint32
		limit     uint32
	)

	if total%req.PageSize > 0 {
		pageTotal += 1
	}

	if pageTotal == 0 {
		pageTotal = 1
	}

	if req.PageIndex > pageTotal {
		req.PageIndex = pageTotal
	}

	offset = (req.PageIndex - 1) * req.PageSize
	limit = req.PageSize

	if offset+limit > total {
		limit = total - offset
	}

	return Pagination{
		PageIndex:   req.PageIndex,
		PageSize:    req.PageSize,
		PageTotal:   pageTotal,
		ResultTotal: total,
		Offset:      int(offset),
		Limit:       int(limit),
	}
}
