// my-article-app/internal/usecase/article_usecase.go
package usecase

import (
	"errors"
	"my-article-app/internal/dto"
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
)

type ArticleUseCase interface {
	CreateArticle(req *dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	GetAllArticles() ([]dto.ArticleResponse, error)
	GetArticleByID(id uint) (*dto.ArticleResponse, error)
	UpdateArticle(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	DeleteArticle(id uint) error
}

type articleUseCase struct {
	articleRepo repository.ArticleRepository
	authorRepo  repository.AuthorRepository // <-- إضافة مستودع المؤلف للتحقق
}

// NewArticleUseCase ينشئ مثيلاً جديدًا (يجب تحديثه في main.go)
func NewArticleUseCase(articleRepo repository.ArticleRepository, authorRepo repository.AuthorRepository) ArticleUseCase {
	return &articleUseCase{
		articleRepo: articleRepo,
		authorRepo:  authorRepo,
	}
}

// CreateArticle ينشئ مقالًا جديدًا
func (uc *articleUseCase) CreateArticle(req *dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	// التحقق من وجود المؤلف أولاً
	author, err := uc.authorRepo.FindByID(req.AuthorID)
	if err != nil || author == nil {
		return nil, errors.New("المؤلف بالمعرف المحدد غير موجود")
	}

	article := &models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: req.AuthorID,
	}

	if err := uc.articleRepo.Create(article); err != nil {
		return nil, err
	}
	
	// GORM لا يقوم بتحميل المؤلف تلقائيًا بعد الإنشاء، لذلك نقوم بتعبئته يدويًا
	article.Author = *author

	return uc.mapArticleToResponse(article), nil
}

// GetAllArticles يجلب جميع المقالات
func (uc *articleUseCase) GetAllArticles() ([]dto.ArticleResponse, error) {
	articles, err := uc.articleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ArticleResponse
	for _, article := range articles {
		responses = append(responses, *uc.mapArticleToResponse(&article))
	}
	return responses, nil
}

// GetArticleByID يجلب مقالًا واحدًا
func (uc *articleUseCase) GetArticleByID(id uint) (*dto.ArticleResponse, error) {
	article, err := uc.articleRepo.FindByID(id)
	if err != nil || article == nil {
		return nil, err
	}
	return uc.mapArticleToResponse(article), nil
}

// UpdateArticle يحدّث مقالًا
func (uc *articleUseCase) UpdateArticle(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	article, err := uc.articleRepo.FindByID(id)
	if err != nil || article == nil {
		return nil, err
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}

	if err := uc.articleRepo.Update(article); err != nil {
		return nil, err
	}
	return uc.mapArticleToResponse(article), nil
}

// DeleteArticle يحذف مقالًا
func (uc *articleUseCase) DeleteArticle(id uint) error {
	return uc.articleRepo.Delete(id)
}

// دالة مساعدة للتحويل لتقليل التكرار
func (uc *articleUseCase) mapArticleToResponse(article *models.Article) *dto.ArticleResponse {
	return &dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		Author: dto.AuthorResponse{
			ID:    article.Author.ID,
			Name:  article.Author.Name,
			Email: article.Author.Email,
		},
	}
}