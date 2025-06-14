// my-article-app/internal/models/article.go
package models

import (
	"time"
)

// Article   بنية قاعدة البيانات فقط
type Article struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	AuthorID  uint
	Author    Author    `gorm:"foreignKey:AuthorID"` // نحتفظ بهذا لـ GORM Preload
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}