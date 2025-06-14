// my-article-app/internal/handlers/author_handler.go
package handlers

import (
	"fmt"
	"log"
	"my-article-app/internal/dto"
	"my-article-app/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


type AuthorHandler interface {
	CreateAuthor(c *fiber.Ctx) error
	GetAllAuthors(c *fiber.Ctx) error
	GetAuthorByID(c *fiber.Ctx) error
	UpdateAuthor(c *fiber.Ctx) error
	DeleteAuthor(c *fiber.Ctx) error
}

type authorHandler struct {
	authorUseCase usecase.AuthorUseCase
}

func NewAuthorHandler(authorUseCase usecase.AuthorUseCase) AuthorHandler {
	return &authorHandler{authorUseCase: authorUseCase}
}

// CreateAuthor يتعامل مع طلبات POST لإنشاء مؤلف جديد
func (h *authorHandler) CreateAuthor(c *fiber.Ctx) error {
	req := new(dto.CreateAuthorRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح."})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "خطأ في التحقق من صحة البيانات.", "details": err.Error()})
	}

	authorResponse, err := h.authorUseCase.CreateAuthor(req)
	if err != nil {
		log.Printf("خطأ في إنشاء المؤلف: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المؤلف."})
	}

	return c.Status(fiber.StatusCreated).JSON(authorResponse)
}

// GetAllAuthors يجلب جميع المؤلفين
func (h *authorHandler) GetAllAuthors(c *fiber.Ctx) error {
	authors, err := h.authorUseCase.GetAllAuthors()
	if err != nil {
		log.Printf("خطأ في جلب المؤلفين: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المؤلفين."})
	}
	return c.JSON(authors)
}

// GetAuthorByID يجلب مؤلفًا واحدًا حسب ID
func (h *authorHandler) GetAuthorByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	author, err := h.authorUseCase.GetAuthorByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المؤلف: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
	}

	return c.JSON(author)
}

// UpdateAuthor يحدّث مؤلفًا موجودًا
func (h *authorHandler) UpdateAuthor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	req := new(dto.UpdateAuthorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح."})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "خطأ في التحقق من صحة البيانات.", "details": err.Error()})
	}

	authorResponse, err := h.authorUseCase.UpdateAuthor(uint(id), req)
	if err != nil {
		log.Printf("خطأ في تحديث المؤلف: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المؤلف."})
	}

	return c.JSON(authorResponse)
}

// DeleteAuthor يحذف مؤلفًا
func (h *authorHandler) DeleteAuthor(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	if err := h.authorUseCase.DeleteAuthor(uint(id)); err != nil {
		log.Printf("خطأ في حذف المؤلف: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود أو فشلت عملية الحذف.", id)})
	}

	return c.SendStatus(fiber.StatusNoContent)
}