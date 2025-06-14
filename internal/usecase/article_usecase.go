// my-article-app/internal/usecase/article_usecase.go
package usecase

import (
	"errors"
	"my-article-app/internal/dto"
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
)

// ArticleUseCase interface remains the same
type ArticleUseCase interface {
	CreateArticle(req *dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	GetAllArticles() ([]dto.ArticleResponse, error)
	GetArticleByID(id uint) (*dto.ArticleResponse, error)
	UpdateArticle(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	DeleteArticle(id uint) error
}

type articleUseCase struct {
	articleRepo repository.ArticleRepository
	authorRepo  repository.AuthorRepository
}

func NewArticleUseCase(articleRepo repository.ArticleRepository, authorRepo repository.AuthorRepository) ArticleUseCase {
	return &articleUseCase{
		articleRepo: articleRepo,
		authorRepo:  authorRepo,
	}
}

// 1. دالة التحويل التي اخترتها (تستقبل المقال والمؤلف بشكل منفصل)
func mapArticleToResponse(article *models.Article, author *models.Author) *dto.ArticleResponse {
	if article == nil || author == nil {
		return nil
	}
	return &dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		Author: dto.AuthorResponse{ // استخدم بيانات المؤلف التي تم تمريرها مباشرة
			ID:    author.ID,
			Name:  author.Name,
			Email: author.Email,
		},
	}
}

// CreateArticle (الحالة الخاصة التي تتطلب جلب المؤلف بشكل منفصل)
func (uc *articleUseCase) CreateArticle(req *dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	// نجلب المؤلف بشكل صريح للتحقق منه
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

	// نمرر المقال الجديد والمؤلف الذي جلبناه إلى دالة التحويل
	return mapArticleToResponse(article, author), nil
}

// GetAllArticles (الحالة العادية)
func (uc *articleUseCase) GetAllArticles() ([]dto.ArticleResponse, error) {
	// Repository's FindAll already preloads the author into each article
	articles, err := uc.articleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ArticleResponse
	for _,article := range articles{


		currentArticle :=article
		responses =append(responses, *mapArticleToResponse(&currentArticle,&currentArticle.Author))
	}
	return responses ,nil
}

// GetArticleByID (الحالة العادية)
func (uc *articleUseCase) GetArticleByID(id uint) (*dto.ArticleResponse, error) {
	// Repository's FindByID already preloads the author
	article, err := uc.articleRepo.FindByID(id)
	if err != nil || article == nil {
		return nil, err
	}

	// نمرر المقال والمؤلف المدمج بداخله (&article.Author) إلى دالة التحويل
	return mapArticleToResponse(article, &article.Author), nil
}

// UpdateArticle (الحالة العادية)
func (uc *articleUseCase) UpdateArticle(id uint, req *dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// Repository's FindByID already preloads the author
	article, err := uc.articleRepo.FindByID(id)
	if err != nil || article == nil {
		return nil, err
	}

	// تحديث الحقول
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}

	if err := uc.articleRepo.Update(article); err != nil {
		return nil, err
	}
	
	// نمرر المقال المحدّث والمؤلف المدمج بداخله (&article.Author) إلى دالة التحويل
	return mapArticleToResponse(article, &article.Author), nil
}

// DeleteArticle remains the same
func (uc *articleUseCase) DeleteArticle(id uint) error {
	return uc.articleRepo.Delete(id)
}
