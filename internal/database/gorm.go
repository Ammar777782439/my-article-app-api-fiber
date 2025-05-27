// my-article-app/internal/database/gorm.go
package database

import (
	"fmt"
	"log"
	"my-article-app/internal/models" // استيراد بنية Article

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitGORMDB تهيئ اتصال GORM بقاعدة بيانات PostgreSQL
func InitGORMDB() (*gorm.DB, error) {
	// هنا نضع سلسلة الاتصال مباشرةً (للتطبيق التعليمي فقط، في الإنتاج استخدم متغيرات البيئة)
	dsn := "host=localhost user=postgres password=postgres dbname=article_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("فشل الاتصال بقاعدة البيانات باستخدام GORM: %w", err)
	}

	// ترحيل (Migrate) مخطط قاعدة البيانات لإنشاء/تحديث جدول Articles
	// AutoMigrate سيقوم بإنشاء الجدول بناءً على بنية Article إذا لم يكن موجودًا.
	// وسيقوم بتحديث الأعمدة إذا أضفت حقولًا جديدة.
	err = db.AutoMigrate(&models.Article{},&models.Author{})
	if err != nil {
		return nil, fmt.Errorf("فشل ترحيل قاعدة البيانات لجدول Articles: %w", err)
	}

	log.Println("تم الاتصال بقاعدة البيانات وترحيل جداول GORM بنجاح!")
	return db, nil
}