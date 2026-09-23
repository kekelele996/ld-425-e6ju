package dto

import "time"

// CreateDesignRequest 创建设计阶段请求。
type CreateDesignRequest struct {
	ProjectID   uint     `json:"project_id" binding:"required"`
	Name        string   `json:"name" binding:"required,max=64"`
	DesignerID  uint     `json:"designer_id"`
	Description string   `json:"description"`
	FileURLs    []string `json:"file_urls"`
}

// UpdateDesignRequest 更新设计阶段请求。
type UpdateDesignRequest struct {
	Name        *string   `json:"name" binding:"omitempty,max=64"`
	DesignerID  *uint     `json:"designer_id"`
	Description *string   `json:"description"`
	FileURLs    *[]string `json:"file_urls"`
}

// SubmitDesignRequest 提交审核请求。
type SubmitDesignRequest struct {
	Description string   `json:"description"`
	FileURLs    []string `json:"file_urls"`
}

// ReviewDesignRequest 审核设计请求。
type ReviewDesignRequest struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment" binding:"max=500"`
}

// DesignDTO 设计阶段展示结构。
type DesignDTO struct {
	ID            uint      `json:"id"`
	ProjectID     uint      `json:"project_id"`
	Name          string    `json:"name"`
	DesignerID    uint      `json:"designer_id"`
	Status        string    `json:"status"`
	Version       int       `json:"version"`
	Description   string    `json:"description"`
	FileURLs      []string  `json:"file_urls"`
	ReviewComment string    `json:"review_comment"`
	ReviewerID    uint      `json:"reviewer_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
