// my-article-app/internal/models/author.go
package models

import (
	"gorm.io/gorm"
)

// Author   بنية قاعدة البيانات فقط
type Author struct {
	gorm.Model
	Name     string
	Email    string    `gorm:"unique;not null"`
	Articles []Article `gorm:"foreignKey:AuthorID"` // نحتفظ بهذا لـ GORM Preload
}