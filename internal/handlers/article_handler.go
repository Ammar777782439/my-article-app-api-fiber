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
	articleRepo *repository.ArticleRepository
	// قد تحتاج هنا لمرجع إلى AuthorRepository إذا كنت بحاجة للتحقق من وجود المؤلف
	// authorRepo  *repository.AuthorRepository // مثال: لو أردت التحقق من AuthorID
}

// NewArticleHandler ينشئ مثيلاً جديدًا من ArticleHandler
func NewArticleHandler(articleRepo *repository.ArticleRepository /*, authorRepo *repository.AuthorRepository*/) *ArticleHandler {
	return &ArticleHandler{articleRepo: articleRepo /*, authorRepo: authorRepo*/}
}

// CreateArticle يتعامل مع طلبات POST لإنشاء مقال جديد
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	article := new(models.Article)
	if err := c.BodyParser(article); err != nil {
		log.Printf("خطأ في تحليل جسم الطلب: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	// التحقق من صحة بيانات المقال (Title, Content, AuthorID)
	// ملاحظة: حقل Author (string) تم إزالته من Article model، الآن تحتاج AuthorID.
	// قد تحتاج إلى إضافة validate:"required" على AuthorID في Article model إذا لم يكن موجوداً.
	if err := validate.Struct(article); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// *** إضافة خطوة اختيارية: التحقق من وجود المؤلف (Author) ***
	// إذا أردت التأكد من أن AuthorID الموجود في الطلب يشير إلى مؤلف حقيقي،
	// ستحتاج إلى جلب AuthorRepository وتعديل NewArticleHandler لقبوله.
	/*
	if article.AuthorID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Author ID مطلوب."})
	}
	author, err := h.authorRepo.FindByID(article.AuthorID)
	if err != nil {
		log.Printf("خطأ في جلب المؤلف من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل التحقق من المؤلف."})
	}
	if author == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", article.AuthorID)})
	}
	*/
	// *** نهاية الخطوة الاختيارية ***

	if err := h.articleRepo.Create(article); err != nil {
		log.Printf("خطأ في إنشاء المقال في قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المقال."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المقال بنجاح!",
		"article": article, // Article الآن سيحتوي على كائن Author المتداخل بفضل Preload في الريبو
	})
}

// GetAllArticles يتعامل مع طلبات GET لجلب جميع المقالات
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	articles, err := h.articleRepo.FindAll()
	if err != nil {
		log.Printf("خطأ في جلب المقالات من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقالات."})
	}
	return c.JSON(articles) // Articles الآن ستتضمن بيانات المؤلفين بسبب Preload
}

// GetArticleByID يتعامل مع طلبات GET لجلب مقال واحد حسب ID
func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	article, err := h.articleRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("خطأ في جلب المقال من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقال."})
	}
	if article == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
	}

	return c.JSON(article) // Article الآن سيتضمن بيانات المؤلف بسبب Preload
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

	// التحقق من صحة البيانات المُرسلة للتحديث
	// لا تنس أن حقل Author (string) لم يعد موجوداً في النموذج
	if err := validate.Struct(article); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	article.ID = uint(id) // تعيين الـ ID للمقال الذي سيتم تحديثه

	// ملاحظة: إذا كان طلب التحديث يرسل AuthorID، فسيتم تحديثه.
	// يمكنك إضافة منطق للتحقق من AuthorID إذا تم توفيره.
	/*
	if article.AuthorID != 0 {
		author, err := h.authorRepo.FindByID(article.AuthorID)
		if err != nil || author == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", article.AuthorID)})
		}
	}
	*/

	if err := h.articleRepo.Update(article); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في تحديث المقال في قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المقال."})
	}

	return c.JSON(fiber.Map{
		"message": "تم تحديث المقال بنجاح!",
		"article": article, // Article الآن سيتضمن بيانات المؤلف بسبب Preload إذا تم تحديثه
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
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		log.Printf("خطأ في حذف المقال من قاعدة البيانات: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المقال."})
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}