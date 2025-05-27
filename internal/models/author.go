// my-article-app/internal/models/author.go
package models

import (
	
	"gorm.io/gorm" // لاستيراد gorm.Model
)


// Author هي بنية لتمثيل المؤلف في التطبيق وقاعدة البيانات.
type Author struct {
  gorm.Model        // يحتوي على ID و CreatedAt و UpdatedAt (بشكل ضمني)
  Name      string  `json:"name" validate:"required,min=3,max=50"`
  Email     string  `json:"email" gorm:"unique;not null" validate:"required,email"`
  
  // تعريف علاقة One-to-Many مع المقالات (اختياري هنا، ولكنه مفيد في نموذج المؤلف)
  // GORM سيقوم تلقائياً بالبحث عن AuthorID في جدول Articles
  Articles  []Article `json:"articles,omitempty" gorm:"foreignKey:AuthorID"`
}