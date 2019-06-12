package model

type Pagination struct {
	Total uint32 `json:"total"`
	Page  int    `json:"page"`
	Count int    `json:"count"`
}

type PaginatedRS struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}
