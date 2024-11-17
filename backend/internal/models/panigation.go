package models

// PaginationQuery represents query parameters for pagination
type PaginationQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	OrderBy  string `form:"order_by,default=id"`
	Sort     string `form:"sort,default=asc"`
}