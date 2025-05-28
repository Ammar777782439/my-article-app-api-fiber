// my-article-app/internal/usecase/author_usecase.go
package usecase

import (
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
)

// AuthorUseCase يمثل منطق العمل للمؤلفين
type AuthorUseCase struct {
	authorRepo *repository.AuthorRepository
}

// NewAuthorUseCase ينشئ مثيلاً جديدًا من AuthorUseCase
func NewAuthorUseCase(authorRepo *repository.AuthorRepository) *AuthorUseCase {
	return &AuthorUseCase{authorRepo: authorRepo}
}

// CreateAuthor يقوم بإنشاء مؤلف جديد
func (uc *AuthorUseCase) CreateAuthor(author *models.Author) error {
	// منطق عمل إضافي (مثال: التحقق من أن البريد الإلكتروني غير مستخدم)
	return uc.authorRepo.Create(author)
}

// GetAllAuthors يجلب جميع المؤلفين
func (uc *AuthorUseCase) GetAllAuthors() ([]models.Author, error) {
	return uc.authorRepo.FindAll()
}

// GetAuthorByID يجلب مؤلفًا واحدًا حسب ID
func (uc *AuthorUseCase) GetAuthorByID(id uint) (*models.Author, error) {
	return uc.authorRepo.FindByID(id)
}

// UpdateAuthor يقوم بتحديث مؤلف موجود
func (uc *AuthorUseCase) UpdateAuthor(author *models.Author) error {
	return uc.authorRepo.Update(author)
}

// DeleteAuthor يقوم بحذف مؤلف
func (uc *AuthorUseCase) DeleteAuthor(id uint) error {
	// منطق عمل إضافي (مثال: التحقق مما إذا كان للمؤلف مقالات قبل الحذف)
	return uc.authorRepo.Delete(id)
}