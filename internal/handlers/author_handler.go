// my-article-app/internal/handlers/author_handler.go
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

// AuthorHandler يمثل معالجات HTTP للمؤلفين
type AuthorHandler struct {
	authorUseCase *usecase.AuthorUseCase // <-- تغيير هنا
}

// NewAuthorHandler ينشئ مثيلاً جديدًا من AuthorHandler
func NewAuthorHandler(authorUseCase *usecase.AuthorUseCase) *AuthorHandler { // <-- تغيير هنا
	return &AuthorHandler{authorUseCase: authorUseCase} // <-- تغيير هنا
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

	if err := h.authorUseCase.CreateAuthor(author); err != nil { // <-- تغيير هنا
		log.Printf("خطأ في إنشاء المؤلف: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المؤلف."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المؤلف بنجاح!",
		"author":  author,
	})
}

// GetAllAuthors يتعامل مع طلبات GET لجلب جميع المؤلفين
func (h *AuthorHandler) GetAllAuthors(c *fiber.Ctx) error {
	authors, err := h.authorUseCase.GetAllAuthors() // <-- تغيير هنا
	if err != nil {
		log.Printf("خطأ في جلب المؤلفين: %v", err)
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

	author, err := h.authorUseCase.GetAuthorByID(uint(id)) // <-- تغيير هنا
	if err != nil {
		log.Printf("خطأ في جلب المؤلف: %v", err)
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

	author.ID = uint(id)

	if err := h.authorUseCase.UpdateAuthor(author); err != nil { // <-- تغيير هنا
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في تحديث المؤلف: %v", err)
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

	if err := h.authorUseCase.DeleteAuthor(uint(id)); err != nil { // <-- تغيير هنا
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في حذف المؤلف: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المؤلف."})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}