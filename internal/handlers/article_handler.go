// my-article-app/internal/handlers/article_handler.go
package handlers

import (
	"fmt"
	"log"
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" // استيراد GORM للتحقق من gorm.ErrRecordNotFound
)

// ... (validator setup remains the same) ...
var validate = validator.New()

// ArticleHandler يمثل معالجات HTTP للمقالات
type ArticleHandler struct {
	articleRepo *repository.ArticleRepository // لاحظ أن النوع هنا أصبح *repository.ArticleRepository
}

// NewArticleHandler ينشئ مثيلاً جديدًا من ArticleHandler
func NewArticleHandler(articleRepo *repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{articleRepo: articleRepo}
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

	if err := h.articleRepo.Create(article); err != nil {
		log.Printf("خطأ في إنشاء المقال في قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المقال."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المقال بنجاح!",
		"article": article,
	})
}

// GetAllArticles يتعامل مع طلبات GET لجلب جميع المقالات
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	articles, err := h.articleRepo.FindAll()
	if err != nil {
		log.Printf("خطأ في جلب المقالات من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقالات."})
	}
	return c.JSON(articles)
}

// GetArticleByID يتعامل مع طلبات GET لجلب مقال واحد حسب ID
func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // استخدام ParseUint لـ uint
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	article, err := h.articleRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المقال من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقال."})
	}
	if article == nil { // إذا لم يتم العثور على المقال
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

	article.ID = uint(id) // تعيين الـ ID للمقال الذي سيتم تحديثه

	if err := h.articleRepo.Update(article); err != nil {
		if err == gorm.ErrRecordNotFound { // التحقق من خطأ GORM
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في تحديث المقال في قاعدة البيانات: %v", err)
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

	if err := h.articleRepo.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound { // التحقق من خطأ GORM
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في حذف المقال من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المقال."})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}