// my-article-app/internal/handlers/article_handler.go
package handlers

import (
	"fmt"
	"log"
	"my-article-app/internal/dto"
	"my-article-app/internal/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// إعادة استخدام نفس المتغير العام
var validate = validator.New()


type ArticleHandler interface {
	CreateArticle(c *fiber.Ctx) error
	GetAllArticles(c *fiber.Ctx) error
	GetArticleByID(c *fiber.Ctx) error
	UpdateArticle(c *fiber.Ctx) error
	DeleteArticle(c *fiber.Ctx) error
}

type articleHandler struct {
	articleUseCase usecase.ArticleUseCase
}

func NewArticleHandler(articleUseCase usecase.ArticleUseCase) ArticleHandler {
	return &articleHandler{articleUseCase: articleUseCase}
}

// CreateArticle يتعامل مع طلبات POST لإنشاء مقال جديد
func (h *articleHandler) CreateArticle(c *fiber.Ctx) error {
	req := new(dto.CreateArticleRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح."})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "خطأ في التحقق من صحة البيانات.", "details": err.Error()})
	}

	articleResponse, err := h.articleUseCase.CreateArticle(req)
	if err != nil {
		log.Printf("خطأ في إنشاء المقال: %v", err)
		// قد يكون الخطأ لأن المؤلف غير موجود
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(articleResponse)
}

// GetAllArticles يجلب جميع المقالات
func (h *articleHandler) GetAllArticles(c *fiber.Ctx) error {
	articles, err := h.articleUseCase.GetAllArticles()
	if err != nil {
		log.Printf("خطأ في جلب المقالات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقالات."})
	}
	return c.JSON(articles)
}

// GetArticleByID يجلب مقالًا واحدًا حسب ID
func (h *articleHandler) GetArticleByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	article, err := h.articleUseCase.GetArticleByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المقال: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
	}

	return c.JSON(article)
}

// UpdateArticle يحدّث مقالًا موجودًا
func (h *articleHandler) UpdateArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	req := new(dto.UpdateArticleRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح."})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "خطأ في التحقق من صحة البيانات.", "details": err.Error()})
	}

	articleResponse, err := h.articleUseCase.UpdateArticle(uint(id), req)
	if err != nil {
		log.Printf("خطأ في تحديث المقال: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المقال."})
	}

	return c.JSON(articleResponse)
}

// DeleteArticle يحذف مقالًا
func (h *articleHandler) DeleteArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	if err := h.articleUseCase.DeleteArticle(uint(id)); err != nil {
		log.Printf("خطأ في حذف المقال: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود أو فشلت عملية الحذف.", id)})
	}

	return c.SendStatus(fiber.StatusNoContent)
}