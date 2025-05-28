// my-article-app/internal/handlers/article_handler.go
package handlers

import (
	"fmt"
	"log"
	"my-article-app/internal/models"
	"my-article-app/internal/usecase" // استيراد UseCase
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

// ArticleHandler يمثل معالجات HTTP للمقالات
type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase 
}

// NewArticleHandler ينشئ مثيلاً جديدًا من ArticleHandler
func NewArticleHandler(articleUseCase *usecase.ArticleUseCase) *ArticleHandler { 
	return &ArticleHandler{articleUseCase: articleUseCase} 
}

// CreateArticle يتعامل مع طلبات POST لإنشاء مقال جديد
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	article := new(models.Article)
	if err := c.BodyParser(article); err != nil {
		log.Printf("خطأ في تحليل جسم الطلب: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	if err := validate.Struct(article); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// استخدام UseCase 
	if err := h.articleUseCase.CreateArticle(article); err != nil { 
		log.Printf("خطأ في إنشاء المقال: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المقال."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المقال بنجاح!",
		"article": article,
	})
}

// GetAllArticles يتعامل مع طلبات GET لجلب جميع المقالات
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	articles, err := h.articleUseCase.GetAllArticles() // <-- تغيير هنا
	if err != nil {
		log.Printf("خطأ في جلب المقالات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقالات."})
	}
	return c.JSON(articles)
}

// GetArticleByID يتعامل مع طلبات GET لجلب مقال واحد حسب ID
func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	article, err := h.articleUseCase.GetArticleByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المقال: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقال."})
	}
	if article == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
	}

	return c.JSON(article)
}

// UpdateArticle يتعامل مع طلبات PUT لتحديث مقال موجود
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	article := new(models.Article)
	if err := c.BodyParser(article); err != nil {
		log.Printf("خطأ في تحليل جسم الطلب للتحديث: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	if err := validate.Struct(article); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	article.ID = uint(id)

	if err := h.articleUseCase.UpdateArticle(article); err != nil { 
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في تحديث المقال: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المقال."})
	}

	return c.JSON(fiber.Map{
		"message": "تم تحديث المقال بنجاح!",
		"article": article,
	})
}

// DeleteArticle يتعامل مع طلبات DELETE لحذف مقال
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	if err := h.articleUseCase.DeleteArticle(uint(id)); err != nil { 
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في حذف المقال: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المقال."})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}