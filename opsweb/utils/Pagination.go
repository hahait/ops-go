package utils

import (
	"gorm.io/gorm"
)

type Pagination struct {
	Page int `form:"page"`
	PageSize int `form:"page_size"`
}

// 通过 db 进行分页
func (p *Pagination) Paginate(db *gorm.DB) *gorm.DB {
	if p.Page > 0 && p.PageSize > 0 {
		return db.Order("id desc").Offset(p.PageSize * (p.Page - 1)).Limit(p.PageSize)
	}
	return db
}

// 校验是否存在 page 和 page_num 要求分页; 以及设置 page 和 page_num 的默认值
func (p *Pagination) CheckPage() ( *Pagination, bool) {
	if p.Page == 0 && p.PageSize == 0 {
		return p, false
	}

	if p.Page == 0 && p.PageSize != 0 {
		p.Page = 1
		return p, true
	}

	if p.Page != 0 && p.PageSize == 0 {
		p.PageSize = 10
		return p, true
	}
	return p, true
}