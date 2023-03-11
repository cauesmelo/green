package data

import (
	"math"

	"github.com/cauesmelo/green/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
	Order        string
}

type Metadata struct {
	CurretPage   int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurretPage:   page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page < 1, "page", "must be greater than zero")
	v.Check(f.Page > math.MaxInt, "page", "page number too big")
	v.Check(f.PageSize < 1, "page_size", "must be greater than zero")
	v.Check(f.PageSize > 50, "page_size", "must be a maximum of 50")
	v.Check(validator.Out(f.Order, "asc", "desc"), "order", "invalid order value")
	v.Check(validator.Out(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
