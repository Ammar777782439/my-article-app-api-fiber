// my-article-app/internal/dto/author_dto.go
package dto

import "time"

// CreateAuthorRequest هو DTO لطلب إنشاء مؤلف جديد
type CreateAuthorRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateAuthorRequest هو DTO لطلب تحديث بيانات المؤلف
type UpdateAuthorRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=50"`
	Email string `json:"email" validate:"omitempty,email"`
}

// AuthorResponse هو DTO القياسي لإرجاع بيانات المؤلف
type AuthorResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AuthorDetailResponse هو DTO لإرجاع بيانات المؤلف مع مقالاته
type AuthorDetailResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	CreatedAt time.Time         `json:"created_at"`
	Articles  []ArticleResponse `json:"articles,omitempty"`
}