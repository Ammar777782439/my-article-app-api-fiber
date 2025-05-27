// my-article-app/internal/repository/author_repository.go
package repository

import (
	"fmt"
	"my-article-app/internal/models"

	"gorm.io/gorm"
)

// AuthorRepository يمثل المستودع الذي يتعامل مع عمليات CRUD للمؤلفين
type AuthorRepository struct {
	db *gorm.DB
}

// NewAuthorRepository ينشئ مثيلاً جديدًا من AuthorRepository
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create ينشئ مؤلفًا جديدًا في قاعدة البيانات
func (r *AuthorRepository) Create(author *models.Author) error {
	result := r.db.Create(author)
	if result.Error != nil {
		return fmt.Errorf("فشل إنشاء المؤلف: %w", result.Error)
	}
	return nil
}

// FindAll يجلب جميع المؤلفين من قاعدة البيانات
func (r *AuthorRepository) FindAll() ([]models.Author, error) {
	var authors []models.Author
	// يمكن استخدام Preload("Articles") هنا إذا أردت جلب المقالات المرتبطة بكل مؤلف
	result := r.db.Find(&authors)
	if result.Error != nil {
		return nil, fmt.Errorf("فشل جلب المؤلفين: %w", result.Error)
	}
	return authors, nil
}

// FindByID يجلب مؤلفًا واحدًا حسب ID
func (r *AuthorRepository) FindByID(id uint) (*models.Author, error) {
	var author models.Author
	// استخدام Preload("Articles") لجلب المقالات المرتبطة بالمؤلف
	result := r.db.Preload("Articles").First(&author, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // لم يتم العثور على المؤلف
		}
		return nil, fmt.Errorf("فشل جلب المؤلف بالمعرف %d: %w", id, result.Error)
	}
	return &author, nil
}

// Update يقوم بتحديث مؤلف موجود في قاعدة البيانات
func (r *AuthorRepository) Update(author *models.Author) error {
	result := r.db.Save(author)
	if result.Error != nil {
		return fmt.Errorf("فشل تحديث المؤلف: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete يقوم بحذف مؤلف من قاعدة البيانات باستخدام ID
func (r *AuthorRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Author{}, id)
	if result.Error != nil {
		return fmt.Errorf("فشل حذف المؤلف: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}