package dtos

// ==================== REQUESTS ====================

type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required,min=2,max=200"`
	Slug         string   `json:"slug" binding:"omitempty,max=200"`
	Description  string   `json:"description" binding:"required"`
	Price        float64  `json:"price" binding:"required,gt=0"`
	ComparePrice *float64 `json:"comparePrice" binding:"omitempty,gt=0"`
	CategoryID   string   `json:"categoryId" binding:"required"`
	Images       []string `json:"images"`
	Sizes        []string `json:"sizes"`
	Colors       []string `json:"colors"`
	Tags         []string `json:"tags"`
	Stock        int      `json:"stock" binding:"gte=0"`
	SKU          *string  `json:"sku"`
	IsFeatured   bool     `json:"isFeatured"`
	IsActive     *bool    `json:"isActive"`
}

type UpdateProductRequest struct {
	Name         *string  `json:"name" binding:"omitempty,min=2,max=200"`
	Slug         *string  `json:"slug"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price" binding:"omitempty,gt=0"`
	ComparePrice *float64 `json:"comparePrice"`
	CategoryID   *string  `json:"categoryId"`
	Images       []string `json:"images"`
	Sizes        []string `json:"sizes"`
	Colors       []string `json:"colors"`
	Tags         []string `json:"tags"`
	Stock        *int     `json:"stock" binding:"omitempty,gte=0"`
	SKU          *string  `json:"sku"`
	IsFeatured   *bool    `json:"isFeatured"`
	IsActive     *bool    `json:"isActive"`
}

// ==================== FILTERS ====================

type ProductFilter struct {
	CategorySlug string  `form:"category"`
	Search       string  `form:"search"`
	MinPrice     float64 `form:"minPrice"`
	MaxPrice     float64 `form:"maxPrice"`
	Size         string  `form:"size"`
	Color        string  `form:"color"`
	Tag          string  `form:"tag"`
	Featured     *bool   `form:"featured"`
	IsActive     *bool   `form:"isActive"`
	Sort         string  `form:"sort"`
	Order        string  `form:"order"`
	Page         int     `form:"page"`
	Limit        int     `form:"limit"`
}

// ==================== RESPONSES ====================

type ProductResponse struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Description  string         `json:"description"`
	Price        float64        `json:"price"`
	ComparePrice *float64       `json:"comparePrice,omitempty"`
	CategoryID   string         `json:"categoryId"`
	Images       []string       `json:"images"`
	Sizes        []string       `json:"sizes"`
	Colors       []string       `json:"colors"`
	Tags         []string       `json:"tags"`
	Stock        int            `json:"stock"`
	SKU          *string        `json:"sku,omitempty"`
	IsFeatured   bool           `json:"isFeatured"`
	IsActive     bool           `json:"isActive"`
	ViewCount    int            `json:"viewCount"`
	SoldCount    int            `json:"soldCount"`
	Category     *CategoryBrief `json:"category,omitempty"`
	Rating       float64        `json:"rating,omitempty"`
	ReviewCount  int            `json:"reviewCount,omitempty"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
}

// ProductListResponse — for public listing (lighter)
type ProductListResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	Price        float64  `json:"price"`
	ComparePrice *float64 `json:"comparePrice,omitempty"`
	Images       []string `json:"images"`
	Stock        int      `json:"stock"`
	IsFeatured   bool     `json:"isFeatured"`
	CategoryID   string   `json:"categoryId"`
}