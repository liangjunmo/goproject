package pagination

type Pagination interface {
	GetPage() uint32
	GetCapacityPerPage() uint32
	GetTotalPages() uint32
	GetTotalRecords() uint32
	GetOffset() int
	GetLimit() int
}

type DefaultPagination struct {
	Page            uint32 `json:"page"`
	CapacityPerPage uint32 `json:"capacity_per_page"`
	TotalPages      uint32 `json:"total_pages"`
	TotalRecords    uint32 `json:"total_records"`
	Offset          int    `json:"offset"`
	Limit           int    `json:"limit"`
}

func (p DefaultPagination) GetPage() uint32 {
	return p.Page
}

func (p DefaultPagination) GetCapacityPerPage() uint32 {
	return p.CapacityPerPage
}

func (p DefaultPagination) GetTotalPages() uint32 {
	return p.TotalPages
}

func (p DefaultPagination) GetTotalRecords() uint32 {
	return p.TotalRecords
}

func (p DefaultPagination) GetOffset() int {
	return p.Offset
}

func (p DefaultPagination) GetLimit() int {
	return p.Limit
}
