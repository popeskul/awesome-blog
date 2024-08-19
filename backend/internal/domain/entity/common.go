package entity

type List[T any] struct {
	Data       []*T        `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

type Response[T any] struct {
	Data       []*T        `json:"data"`
	Pagination *Pagination `json:"pagination"`
}
