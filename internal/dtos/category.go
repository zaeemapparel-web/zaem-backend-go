package dtos

// ==================== REQUESTS ====================

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Slug        string  `json:"slug" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	Image       *string `json:"image" binding:"omitempty,url"`
	ParentID    *string `json:"parentId"`
	SortOrder   int     `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Slug        *string `json:"slug" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	ParentID    *string `json:"parentId"`
	SortOrder   *int    `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

// ==================== RESPONSES ====================

type CategoryResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description *string            `json:"description,omitempty"`
	Image       *string            `json:"image,omitempty"`
	ParentID    *string            `json:"parentId,omitempty"`
	IsParent    bool               `json:"isParent"`
	SortOrder   int                `json:"sortOrder"`
	IsActive    bool               `json:"isActive"`
	Children    []CategoryResponse `json:"children,omitempty"`
	Parent      *CategoryBrief     `json:"parent,omitempty"`
	ProductCount int64             `json:"productCount,omitempty"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

type CategoryBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}