// my-article-app/internal/usecase/author_usecase.go
package usecase

import (
	"my-article-app/internal/dto"
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
)

type AuthorUseCase interface {
	CreateAuthor(req *dto.CreateAuthorRequest) (*dto.AuthorResponse, error)
	GetAllAuthors() ([]dto.AuthorResponse, error)
	GetAuthorByID(id uint) (*dto.AuthorDetailResponse, error)
	UpdateAuthor(id uint, req *dto.UpdateAuthorRequest) (*dto.AuthorResponse, error)
	DeleteAuthor(id uint) error
}

type authorUseCase struct {
	authorRepo repository.AuthorRepository
}

func NewAuthorUseCase(authorRepo repository.AuthorRepository) AuthorUseCase {
	return &authorUseCase{authorRepo: authorRepo}
}

// CreateAuthor ينشئ مؤلفًا جديدًا
func (uc *authorUseCase) CreateAuthor(req *dto.CreateAuthorRequest) (*dto.AuthorResponse, error) {
	author := &models.Author{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := uc.authorRepo.Create(author); err != nil {
		return nil, err
	}

	response := &dto.AuthorResponse{
		ID:    author.ID,
		Name:  author.Name,
		Email: author.Email,
	}
	return response, nil
}

// GetAllAuthors يجلب جميع المؤلفين
func (uc *authorUseCase) GetAllAuthors() ([]dto.AuthorResponse, error) {
	authors, err := uc.authorRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.AuthorResponse
	for _, author := range authors {
		responses = append(responses, dto.AuthorResponse{
			ID:    author.ID,
			Name:  author.Name,
			Email: author.Email,
		})
	}
	return responses, nil
}

// GetAuthorByID يجلب مؤلفًا واحدًا مع مقالاته
func (uc *authorUseCase) GetAuthorByID(id uint) (*dto.AuthorDetailResponse, error) {
	author, err := uc.authorRepo.FindByID(id)
	if err != nil || author == nil {
		return nil, err
	}

	response := &dto.AuthorDetailResponse{
		ID:        author.ID,
		Name:      author.Name,
		Email:     author.Email,
		CreatedAt: author.CreatedAt,
		Articles:  []dto.ArticleResponse{}, // Initialize to avoid null
	}

	for _, article := range author.Articles {
		response.Articles = append(response.Articles, dto.ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
			// Note: Author data is omitted here to avoid circular nesting
		})
	}

	return response, nil
}

// UpdateAuthor يحدّث بيانات المؤلف
func (uc *authorUseCase) UpdateAuthor(id uint, req *dto.UpdateAuthorRequest) (*dto.AuthorResponse, error) {
	author, err := uc.authorRepo.FindByID(id)
	if err != nil || author == nil {
		return nil, err
	}

	if req.Name != "" {
		author.Name = req.Name
	}
	if req.Email != "" {
		author.Email = req.Email
	}

	if err := uc.authorRepo.Update(author); err != nil {
		return nil, err
	}

	response := &dto.AuthorResponse{
		ID:    author.ID,
		Name:  author.Name,
		Email: author.Email,
	}
	return response, nil
}

// DeleteAuthor يحذف المؤلف
func (uc *authorUseCase) DeleteAuthor(id uint) error {
	// Optional: Add logic here to check if the author has articles before deleting.
	return uc.authorRepo.Delete(id)
}