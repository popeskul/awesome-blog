package entity

import (
	"fmt"
)

const (
	MaxLimit     = 100
	DefaultLimit = 10
	MinPage      = 1
)

type RemoteParams struct {
	Page   *int    `form:"page,omitempty" json:"page,omitempty"`
	Limit  *int    `form:"limit,omitempty" json:"limit,omitempty"`
	Offset *int    `form:"offset,omitempty" json:"offset,omitempty"`
	Sort   *string `form:"sort,omitempty" json:"sort,omitempty"`
}

type Pagination struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	Total  int    `json:"total"`
}

func NewPaginationFromParams(params RemoteParams) (*Pagination, error) {
	page := 1
	limit := 10
	offset := 0
	sort := "created_at_desc"

	if params.Page != nil {
		page = *params.Page
	}

	if params.Limit != nil {
		limit = *params.Limit
	}

	if params.Offset != nil {
		offset = *params.Offset
	}

	if params.Sort != nil && *params.Sort != "" {
		switch *params.Sort {
		case "created_at_asc", "created_at_desc", "title_asc", "title_desc":
			sort = *params.Sort
		default:
			return nil, fmt.Errorf("invalid sort parameter: %s", *params.Sort)
		}
	}

	pagination := &Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	fmt.Println("Pagination", pagination)

	if err := ValidatePagination(pagination); err != nil {
		return nil, err
	}

	return pagination, nil
}

func ValidatePagination(params *Pagination) error {
	if params == nil {
		return fmt.Errorf("pagination params cannot be nil")
	}

	if params.Page < MinPage {
		params.Page = MinPage
	}

	if params.Limit <= 0 {
		params.Limit = DefaultLimit
	} else if params.Limit > MaxLimit {
		return fmt.Errorf("limit must be less than or equal to %d", MaxLimit)
	}

	if params.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}

	if params.Offset == 0 && params.Page > MinPage {
		params.Offset = (params.Page - 1) * params.Limit
	}

	if params.Offset > 0 && params.Page == MinPage {
		params.Page = (params.Offset / params.Limit) + 1
	}

	return nil
}
