// my-article-app/internal/dto/article_dto.go
package dto

import "time"

// CreateArticleRequest هو DTO لطلب إنشاء مقال جديد
type CreateArticleRequest struct {
	Title    string `json:"title" validate:"required,min=5,max=200"`
	Content  string `json:"content" validate:"required,min=10"`
	AuthorID uint   `json:"author_id" validate:"required"`
}

// UpdateArticleRequest هو DTO لطلب تحديث مقال
type UpdateArticleRequest struct {
	Title   string `json:"title" validate:"omitempty,min=5,max=200"`
	Content string `json:"content" validate:"omitempty,min=10"`
}

// ArticleResponse هو DTO لإرجاع بيانات المقال
type ArticleResponse struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Author    AuthorResponse `json:"author"`
}