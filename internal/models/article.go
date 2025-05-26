// my-article-app/internal/models/article.go
package models

import "time" // لاستخدام time.Time لحقول التواريخ

// Article هي بنية لتمثيل المقال في التطبيق وقاعدة البيانات.
type Article struct {
	// GORM سيقوم تلقائياً بتعيين ID كـ PRIMARY KEY تزايدي.
	ID        uint      `json:"id" gorm:"primaryKey"` // gorm:"primaryKey" يحدد المفتاح الأساسي
	Title     string    `json:"title" validate:"required,min=5,max=200"`
	Content   string    `json:"content" validate:"required,min=10"`
	Author    string    `json:"author" validate:"required,min=3,max=50"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // GORM سيقوم بتعيين هذا تلقائياً عند الإنشاء
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // GORM سيقوم بتحديث هذا تلقائياً عند التحديث
}