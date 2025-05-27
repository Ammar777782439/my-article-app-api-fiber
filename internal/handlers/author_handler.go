// my-article-app/internal/handlers/author_handler.go
package handlers

import (
	"fmt"
	"log"
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// validator هنا هو نفسه الذي تم تعريفه في article_handler.go
// من الأفضل عادةً تعريف validator مرة واحدة في مكان مركزي إذا كان سيستخدم في عدة handlers.
// للحفاظ على البساطة هنا، نتركه كما هو.
// var validate = validator.New() // إذا كان بالفعل معرفًا في ملف آخر بنفس الحزمة، لا تعيد تعريفه

// AuthorHandler يمثل معالجات HTTP للمؤلفين
type AuthorHandler struct {
	authorRepo *repository.AuthorRepository
}

// NewAuthorHandler ينشئ مثيلاً جديدًا من AuthorHandler
func NewAuthorHandler(authorRepo *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{authorRepo: authorRepo}
}

// CreateAuthor يتعامل مع طلبات POST لإنشاء مؤلف جديد
func (h *AuthorHandler) CreateAuthor(c *fiber.Ctx) error {
	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		log.Printf("خطأ في تحليل جسم طلب المؤلف: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	if err := validate.Struct(author); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	if err := h.authorRepo.Create(author); err != nil {
		log.Printf("خطأ في إنشاء المؤلف في قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المؤلف."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المؤلف بنجاح!",
		"author":  author,
	})
}

// GetAllAuthors يتعامل مع طلبات GET لجلب جميع المؤلفين
func (h *AuthorHandler) GetAllAuthors(c *fiber.Ctx) error {
	authors, err := h.authorRepo.FindAll()
	if err != nil {
		log.Printf("خطأ في جلب المؤلفين من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المؤلفين."})
	}
	return c.JSON(authors)
}

// GetAuthorByID يتعامل مع طلبات GET لجلب مؤلف واحد حسب ID
func (h *AuthorHandler) GetAuthorByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	author, err := h.authorRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المؤلف من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المؤلف."})
	}
	if author == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
	}

	return c.JSON(author)
}

// UpdateAuthor يتعامل مع طلبات PUT لتحديث مؤلف موجود
func (h *AuthorHandler) UpdateAuthor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		log.Printf("خطأ في تحليل جسم طلب التحديث للمؤلف: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	if err := validate.Struct(author); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	author.ID = uint(id) // تعيين الـ ID للمؤلف الذي سيتم تحديثه

	if err := h.authorRepo.Update(author); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في تحديث المؤلف في قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المؤلف."})
	}

	return c.JSON(fiber.Map{
		"message": "تم تحديث المؤلف بنجاح!",
		"author":  author,
	})
}

// DeleteAuthor يتعامل مع طلبات DELETE لحذف مؤلف
func (h *AuthorHandler) DeleteAuthor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	if err := h.authorRepo.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في حذف المؤلف من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المؤلف."})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}