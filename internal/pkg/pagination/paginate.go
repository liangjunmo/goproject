package pagination

func Paginate(page, defaultPage, capacity, defaultCapacity, totalRecords uint32) Pagination {
	if page < 1 {
		page = defaultPage
	}

	if capacity < 1 {
		capacity = defaultCapacity
	}

	var (
		totalPages = totalRecords / capacity
		offset     uint32
		limit      uint32
	)

	if totalRecords%capacity > 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	offset = (page - 1) * capacity
	limit = capacity

	if offset+limit > totalRecords {
		limit = totalRecords - offset
	}

	return DefaultPagination{
		Page:            page,
		CapacityPerPage: capacity,
		TotalPages:      totalPages,
		TotalRecords:    totalRecords,
		Offset:          int(offset),
		Limit:           int(limit),
	}
}
