// my-article-app/internal/models/article.go
package models

import (
	"time" 
)
// sliceBigInt
// Article هي بنية لتمثيل المقال في التطبيق وقاعدة البيانات.utls  /requst /response /role /assert /help 
type Article struct {
	ID      uint   `json:"id" gorm:"primaryKey"` // gorm:"primaryKey" يحدد المفتاح الأساسي
	Title   string `json:"title" validate:"required,min=5,max=200"`
	Content string `json:"content" validate:"required,min=10"`
	// Author    string    `json:"author" validate:"required,min=3,max=50"` // تم إزالة هذا الحقل

	AuthorID uint   `json:"author_id" validate:"required"`                            // تأكد من أن AuthorID مطلوب أيضاً
	Author   Author `json:"author,omitempty" gorm:"foreignKey:AuthorID" validate:"-"` // <-- أضف validate:"-" هنا

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // GORM سيقوم بتعيين هذا تلقائياً عند الإنشاء
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // GORM سيقوم بتحديث هذا تلقائياً عند التحديث
}


// type Author struct {
//   gorm.Model        // يحتوي على ID و CreatedAt و UpdatedAt (بشكل ضمني)
//   Name      string  `json:"name" validate:"required,min=3,max=50"` // أضفت validate tags
//   Email     string  `json:"email" gorm:"unique;not null" validate:"required,email"` // تم تصحيح 'email' إلى 'Email' وأضفت validate tags


//   Articles  []Article `json:"articles,omitempty" gorm:"foreignKey:AuthorID"`
// }
