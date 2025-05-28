// my-article-app/internal/usecase/article_usecase.go
package usecase

import (
	"my-article-app/internal/models"
	"my-article-app/internal/repository"

	

)


// ArticleUseCase يمثل منطق العمل للمقالات
type ArticleUseCase struct {
	articleRepo *repository.ArticleRepository
	// يمكنك إضافة AuthorRepository هنا إذا احتجت للتحقق من المؤلف
	// authorRepo  *repository.AuthorRepository
}

// NewArticleUseCase ينشئ مثيلاً جديدًا من ArticleUseCase
func NewArticleUseCase(articleRepo *repository.ArticleRepository /*, authorRepo *repository.AuthorRepository*/) *ArticleUseCase {
	return &ArticleUseCase{articleRepo: articleRepo /*, authorRepo: authorRepo*/}
}

// CreateArticle يقوم بإنشاء مقال جديد (هنا يمكن إضافة منطق عمل إضافي)
func (uc *ArticleUseCase) CreateArticle(article *models.Article) error {
	// هنا يمكن إضافة أي منطق عمل قبل الحفظ
	// مثال: التحقق من وجود المؤلف (إذا تم تمرير authorRepo)
	// author, err := uc.authorRepo.FindByID(article.AuthorID)
	// if err != nil || author == nil {
	// 	return errors.New("المؤلف غير موجود")
	// }
	
	

	return uc.articleRepo.Create(article)
}

// GetAllArticles يجلب جميع المقالات
func (uc *ArticleUseCase) GetAllArticles() ([]models.Article, error) {
	return uc.articleRepo.FindAll()
}

// GetArticleByID يجلب مقالًا واحدًا حسب ID
func (uc *ArticleUseCase) GetArticleByID(id uint) (*models.Article, error) {
	return uc.articleRepo.FindByID(id)
}

// UpdateArticle يقوم بتحديث مقال موجود
func (uc *ArticleUseCase) UpdateArticle(article *models.Article) error {
	// هنا يمكن إضافة منطق عمل قبل التحديث
	return uc.articleRepo.Update(article)
}

// DeleteArticle يقوم بحذف مقال
func (uc *ArticleUseCase) DeleteArticle(id uint) error {
	// هنا يمكن إضافة منطق عمل قبل الحذف (مثل التحقق من الأذونات)
	return uc.articleRepo.Delete(id)
}