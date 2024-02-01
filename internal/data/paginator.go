package data

import "finance/internal/validator"

type Paginator struct {
	Page     int
	PageSize int
}

func (f Paginator) limit() int {
	return f.PageSize
}

func (f Paginator) offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *validator.Validator, p Paginator) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")
}

