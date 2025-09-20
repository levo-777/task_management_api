package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination metadata for responses
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// FilterParams represents filtering parameters
type FilterParams struct {
	Search     string            `json:"search"`
	SortBy     string            `json:"sort_by"`
	SortOrder  string            `json:"sort_order"`
	Filters    map[string]string `json:"filters"`
}

// GetPaginationParams extracts pagination parameters from Gin context
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// GetFilterParams extracts filtering parameters from Gin context
func GetFilterParams(c *gin.Context) FilterParams {
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Extract additional filters
	filters := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 && key != "page" && key != "page_size" && key != "search" && key != "sort_by" && key != "sort_order" {
			filters[key] = values[0]
		}
	}

	return FilterParams{
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters:   filters,
	}
}

// CreatePaginationResponse creates a paginated response
func CreatePaginationResponse(data interface{}, total int64, params PaginationParams) PaginationResponse {
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return PaginationResponse{
		Data: data,
		Pagination: Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}
}

// ApplySorting applies sorting to a GORM query
func ApplySorting(db *gorm.DB, sortBy, sortOrder string, allowedFields []string) *gorm.DB {
	// Check if sortBy is in allowed fields
	allowed := false
	for _, field := range allowedFields {
		if field == sortBy {
			allowed = true
			break
		}
	}

	if !allowed {
		sortBy = "created_at"
	}

	return db.Order(sortBy + " " + sortOrder)
}

// ApplySearch applies search functionality to a GORM query
func ApplySearch(db *gorm.DB, search string, searchFields []string) *gorm.DB {
	if search == "" || len(searchFields) == 0 {
		return db
	}

	searchPattern := "%" + search + "%"
	
	// Create OR conditions for each search field
	query := ""
	args := make([]interface{}, len(searchFields)*2)
	
	for i, field := range searchFields {
		if i > 0 {
			query += " OR "
		}
		query += field + " ILIKE ?"
		args[i] = searchPattern
	}

	return db.Where(query, args...)
}

// ApplyFilters applies additional filters to a GORM query
func ApplyFilters(db *gorm.DB, filters map[string]string, allowedFilters []string) *gorm.DB {
	for key, value := range filters {
		// Check if filter is allowed
		allowed := false
		for _, field := range allowedFilters {
			if field == key {
				allowed = true
				break
			}
		}

		if allowed && value != "" {
			db = db.Where(key+" = ?", value)
		}
	}

	return db
}
