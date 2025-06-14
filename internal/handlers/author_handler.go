// my-article-app/internal/handlers/author_handler.go
package handlers

// استيراد المكتبات اللازمة للعمل
import (
	"fmt"                          
	"log"                           
	"my-article-app/internal/models"  
	"my-article-app/internal/usecase" 
	"strconv"                       

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
type AuthorHandler interface {
	// تعريف واجهة AuthorHandler التي تحتوي على دوال لمعالجة طلبات HTTP للمؤلفين
	CreateAuthor(c *fiber.Ctx) error // دالة لإنشاء مؤلف جديد
	GetAllAuthors(c *fiber.Ctx) error // دالة لجلب جميع المؤلفين
	GetAuthorByID(c *fiber.Ctx) error // دالة لجلب مؤلف واحد حسب المعرف
	UpdateAuthor(c *fiber.Ctx) error // دالة لتحديث مؤلف موجود
	DeleteAuthor(c *fiber.Ctx) error // دالة لحذف مؤلف موجود
}

// AuthorHandler يمثل معالجات HTTP للمؤلفين

type authorHandler struct {
	authorUseCase usecase.AuthorUseCase // مرجع لطبقة منطق العمل للمؤلفين (الطبقة الثانية)
}

// NewAuthorHandler ينشئ مثيلاً جديدًا من AuthorHandler
// هذه الدالة تُستخدم عند بدء التطبيق لربط طبقة المعالج (Handler) بطبقة منطق العمل (UseCase)
func NewAuthorHandler(authorUseCase usecase.AuthorUseCase) AuthorHandler {
	// إنشاء وإرجاع كائن المعالج مع تمرير مرجع لطبقة منطق العمل
	return &authorHandler{authorUseCase: authorUseCase}
	
}

// CreateAuthor يتعامل مع طلبات POST لإنشاء مؤلف جديد
// هذه الدالة تستقبل طلب HTTP من المستخدم لإنشاء مؤلف جديد
func (h *authorHandler) CreateAuthor(c *fiber.Ctx) error {
	// 1. إنشاء كائن مؤلف جديد فارغ لاستقبال البيانات
	author := new(models.Author)

	// 2. تحليل جسم الطلب (JSON) وتحويله إلى كائن Author
	if err := c.BodyParser(author); err != nil {
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في تحليل جسم طلب المؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم بتنسيق JSON
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	// 3. التحقق من صحة البيانات باستخدام مكتبة validator
	if err := validate.Struct(author); err != nil {
		// تحويل الخطأ إلى نوع ValidationErrors للحصول على تفاصيل أكثر
		validationErrors := err.(validator.ValidationErrors)
		// إرجاع رسالة خطأ مع تفاصيل المشكلة
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// 4. استدعاء طبقة منطق العمل (UseCase) لإنشاء المؤلف في قاعدة البيانات
	if err := h.authorUseCase.CreateAuthor(author); err != nil {
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في إنشاء المؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المؤلف."})
	}

	// 5. إرجاع استجابة ناجحة مع المؤلف الذي تم إنشاؤه
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المؤلف بنجاح!",
		"author":  author, // إرجاع بيانات المؤلف بعد إنشائه (بما في ذلك الـ ID المولد)
	})
}

// GetAllAuthors يتعامل مع طلبات GET لجلب جميع المؤلفين
// هذه الدالة تستقبل طلب HTTP من المستخدم لعرض قائمة بجميع المؤلفين
func (h *authorHandler) GetAllAuthors(c *fiber.Ctx) error {
	// 1. استدعاء طبقة منطق العمل (UseCase) لجلب جميع المؤلفين من قاعدة البيانات
	authors, err := h.authorUseCase.GetAllAuthors()

	// 2. التحقق من وجود أخطاء أثناء جلب البيانات
	if err != nil {
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في جلب المؤلفين: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المؤلفين."})
	}

	// 3. إرجاع قائمة المؤلفين كاستجابة JSON
	return c.JSON(authors)
}

// GetAuthorByID يتعامل مع طلبات GET لجلب مؤلف واحد حسب ID
// هذه الدالة تستقبل طلب HTTP من المستخدم لعرض مؤلف محدد بواسطة المعرف الخاص به
func (h *authorHandler) GetAuthorByID(c *fiber.Ctx) error {
	// 1. استخراج معرف المؤلف من معلمات الطلب (URL parameters)
	idStr := c.Params("id")

	// 2. تحويل المعرف من نص إلى رقم
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إرجاع رسالة خطأ إذا كان المعرف غير صالح
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	// 3. استدعاء طبقة منطق العمل (UseCase) لجلب المؤلف من قاعدة البيانات
	author, err := h.authorUseCase.GetAuthorByID(uint(id))

	// 4. التحقق من وجود أخطاء أثناء جلب البيانات
	if err != nil {
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في جلب المؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المؤلف."})
	}

	// 5. التحقق من وجود المؤلف
	if author == nil {
		// إرجاع رسالة خطأ 404 إذا لم يتم العثور على المؤلف
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
	}

	// 6. إرجاع بيانات المؤلف كاستجابة JSON
	return c.JSON(author)
}

// UpdateAuthor يتعامل مع طلبات PUT لتحديث مؤلف موجود
// هذه الدالة تستقبل طلب HTTP من المستخدم لتحديث بيانات مؤلف موجود
func (h *authorHandler) UpdateAuthor(c *fiber.Ctx) error {
	// 1. استخراج معرف المؤلف من معلمات الطلب (URL parameters)
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إرجاع رسالة خطأ إذا كان المعرف غير صالح
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	// 2. إنشاء كائن مؤلف جديد لاستقبال البيانات المحدثة
	author := new(models.Author)

	// 3. تحليل جسم الطلب (JSON) وتحويله إلى كائن Author
	if err := c.BodyParser(author); err != nil {
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في تحليل جسم طلب التحديث للمؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	// 4. التحقق من صحة البيانات باستخدام مكتبة validator
	if err := validate.Struct(author); err != nil {
		// تحويل الخطأ إلى نوع ValidationErrors للحصول على تفاصيل أكثر
		validationErrors := err.(validator.ValidationErrors)
		// إرجاع رسالة خطأ مع تفاصيل المشكلة
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// 5. تعيين معرف المؤلف من المعلمات لضمان تحديث المؤلف الصحيح
	author.ID = uint(id)

	// 6. استدعاء طبقة منطق العمل (UseCase) لتحديث المؤلف في قاعدة البيانات
	if err := h.authorUseCase.UpdateAuthor(author); err != nil {
		// التحقق من نوع الخطأ - إذا كان المؤلف غير موجود
		if err == gorm.ErrRecordNotFound {
			// إرجاع رسالة خطأ 404 إذا لم يتم العثور على المؤلف
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في تحديث المؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المؤلف."})
	}

	// 7. إرجاع استجابة ناجحة مع المؤلف المحدث
	return c.JSON(fiber.Map{
		"message": "تم تحديث المؤلف بنجاح!",
		"author":  author,
	})
}

// DeleteAuthor يتعامل مع طلبات DELETE لحذف مؤلف
// هذه الدالة تستقبل طلب HTTP من المستخدم لحذف مؤلف موجود
func (h *authorHandler) DeleteAuthor(c *fiber.Ctx) error {
	// 1. استخراج معرف المؤلف من معلمات الطلب (URL parameters)
	idStr := c.Params("id")

	// 2. تحويل المعرف من نص إلى رقم
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إرجاع رسالة خطأ إذا كان المعرف غير صالح
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المؤلف غير صالح."})
	}

	// 3. استدعاء طبقة منطق العمل (UseCase) لحذف المؤلف من قاعدة البيانات
	if err := h.authorUseCase.DeleteAuthor(uint(id)); err != nil {
		// التحقق من نوع الخطأ - إذا كان المؤلف غير موجود
		if err == gorm.ErrRecordNotFound {
			// إرجاع رسالة خطأ 404 إذا لم يتم العثور على المؤلف
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المؤلف بالمعرف %d غير موجود.", id)})
		}
		// تسجيل الخطأ للمتابعة وتصحيح المشكلات
		log.Printf("خطأ في حذف المؤلف: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المؤلف."})
	}

	// 4. إرجاع استجابة ناجحة بدون محتوى (204 No Content)
	// هذا هو السلوك القياسي لعمليات الحذف الناجحة في REST API
	return c.Status(fiber.StatusNoContent).SendString("")
}