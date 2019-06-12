package util

import "bitbucket.org/inehealth/idonia-common/filter"

//Ascendant sql string
const Ascendant = "ASC"

//Descendant sql string
const Descendant = "DESC"

//PageResult pagination struct
type PageResult struct {
	CurrentPage        int                 `json:"current_page"`
	TotalPages         int                 `json:"total_pages"`
	ResultPerPageCount int                 `json:"result_per_page_count"`
	ResultTotalCount   int                 `json:"result_total_count"`
	Orders             []Order             `json:"order"`
	Filter             filter.SearchFilter `json:"filter"`
}

//Order pagination struct
type Order struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}
