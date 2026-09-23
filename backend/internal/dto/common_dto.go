package dto

// PageQuery 通用分页参数。
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// Normalize 返回规范化分页值。
func (q PageQuery) Normalize() (page, pageSize int) {
	page = q.Page
	if page <= 0 {
		page = 1
	}
	pageSize = q.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	return page, pageSize
}
