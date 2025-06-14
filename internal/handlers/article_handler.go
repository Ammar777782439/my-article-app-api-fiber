// my-article-app/internal/handlers/article_handler.go
package handlers


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

// إنشاء متغير عام للتحقق من صحة البيانات
var validate = validator.New()

type ArticleHandler interface {
	// تعريف واجهة ArticleHandler التي تحتوي على دوال لمعالجة طلبات HTTP للمقالات
	CreateArticle(c *fiber.Ctx) error 
	GetAllArticles(c *fiber.Ctx) error 
	GetArticleByID(c *fiber.Ctx) error
	UpdateArticle(c *fiber.Ctx) error 
	DeleteArticle(c *fiber.Ctx) error
}
// ArticleHandler يمثل معالجات HTTP للمقالات (هذا هو الطابق الأول الذي يتعامل مع طلبات المستخدم)

type articleHandler struct {
	articleUseCase usecase.ArticleUseCase // مرجع لطبقة منطق العمل للمقالات (الطابق الثاني)
}

// NewArticleHandler ينشئ مثيلاً جديدًا من ArticleHandler
// هذه الدالة تُستخدم عند بدء التطبيق لربط طبقة المعالج (Handler) بطبقة منطق العمل (UseCase)
func NewArticleHandler(articleUseCase usecase.ArticleUseCase) ArticleHandler { 
	// إنشاء وإرجاع معالج جديد مع تمرير مرجع لطبقة منطق العمل
	return &articleHandler{articleUseCase: articleUseCase} 
}

// CreateArticle يتعامل مع طلبات POST لإنشاء مقال جديد
// عندما يريد المستخدم إنشاء مقال جديد، يتم استدعاء هذه الدالة
func (h *articleHandler) CreateArticle(c *fiber.Ctx) error {
	// 1. إنشاء مقال جديد فارغ لاستقبال البيانات
	article := new(models.Article)

	// 2. استخراج بيانات المقال من جسم الطلب (JSON)
	if err := c.BodyParser(article); err != nil {
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في تحليل جسم الطلب: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	// 3. التحقق من صحة البيانات (مثل: هل العنوان موجود؟ هل المحتوى غير فارغ؟)
	if err := validate.Struct(article); err != nil {
		// تحويل أخطاء التحقق إلى تنسيق مفهوم
		validationErrors := err.(validator.ValidationErrors)
		// إرجاع تفاصيل الأخطاء للمستخدم
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// 4. استدعاء طبقة منطق العمل (الطابق الثاني) لإنشاء المقال
	if err := h.articleUseCase.CreateArticle(article); err != nil { 
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في إنشاء المقال: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل إنشاء المقال."})
	}

	// 5. إذا نجحت العملية، إرسال رد النجاح للمستخدم مع بيانات المقال المنشأ
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "تم إنشاء المقال بنجاح!",
		"article": article,
	})
}

// GetAllArticles يتعامل مع طلبات GET لجلب جميع المقالات
// عندما يريد المستخدم عرض قائمة بجميع المقالات
func (h *articleHandler) GetAllArticles(c *fiber.Ctx) error {
	// 1. استدعاء طبقة منطق العمل (الطابق الثاني) لجلب جميع المقالات
	articles, err := h.articleUseCase.GetAllArticles() 

	// 2. التحقق من وجود أخطاء
	if err != nil {
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في جلب المقالات: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقالات."})
	}

	// 3. إرجاع قائمة المقالات للمستخدم بتنسيق JSON
	return c.JSON(articles)
}

// GetArticleByID يتعامل مع طلبات GET لجلب مقال واحد حسب ID
// عندما يريد المستخدم عرض مقال محدد باستخدام رقم المعرف
func (h *articleHandler) GetArticleByID(c *fiber.Ctx) error {
	// 1. استخراج معرف المقال من عنوان URL
	idStr := c.Params("id")

	// 2. تحويل المعرف من نص إلى رقم
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إذا كان المعرف غير صالح (ليس رقماً)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	// 3. استدعاء طبقة منطق العمل (الطابق الثاني) لجلب المقال
	article, err := h.articleUseCase.GetArticleByID(uint(id))

	// 4. التحقق من وجود أخطاء
	if err != nil {
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في جلب المقال: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل جلب المقال."})
	}

	// 5. التحقق مما إذا كان المقال موجوداً
	if article == nil {
		// إذا لم يتم العثور على المقال
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
	}

	// 6. إرجاع المقال للمستخدم بتنسيق JSON
	return c.JSON(article)
}

// UpdateArticle يتعامل مع طلبات PUT لتحديث مقال موجود
// عندما يريد المستخدم تحديث بيانات مقال موجود
func (h *articleHandler) UpdateArticle(c *fiber.Ctx) error {
	// 1. استخراج معرف المقال من عنوان URL
	idStr := c.Params("id")

	// 2. تحويل المعرف من نص إلى رقم
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إذا كان المعرف غير صالح (ليس رقماً)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	// 3. إنشاء مقال جديد فارغ لاستقبال البيانات المحدثة
	article := new(models.Article)

	// 4. استخراج بيانات المقال المحدثة من جسم الطلب (JSON)
	if err := c.BodyParser(article); err != nil {
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في تحليل جسم الطلب للتحديث: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "جسم الطلب غير صالح أو مفقود."})
	}

	// 5. التحقق من صحة البيانات المحدثة
	if err := validate.Struct(article); err != nil {
		// تحويل أخطاء التحقق إلى تنسيق مفهوم
		validationErrors := err.(validator.ValidationErrors)
		// إرجاع تفاصيل الأخطاء للمستخدم
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "خطأ في التحقق من صحة البيانات.",
			"details": validationErrors.Error(),
		})
	}

	// 6. تعيين معرف المقال الذي نريد تحديثه
	article.ID = uint(id)

	// 7. استدعاء طبقة منطق العمل (الطابق الثاني) لتحديث المقال
	if err := h.articleUseCase.UpdateArticle(article); err != nil { 
		// إذا لم يتم العثور على المقال
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في تحديث المقال: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل تحديث المقال."})
	}

	// 8. إذا نجحت العملية، إرسال رد النجاح للمستخدم مع بيانات المقال المحدث
	return c.JSON(fiber.Map{
		"message": "تم تحديث المقال بنجاح!",
		"article": article,
	})
}

// DeleteArticle يتعامل مع طلبات DELETE لحذف مقال
// عندما يريد المستخدم حذف مقال موجود
func (h *articleHandler) DeleteArticle(c *fiber.Ctx) error {
	// 1. استخراج معرف المقال من عنوان URL
	idStr := c.Params("id")

	// 2. تحويل المعرف من نص إلى رقم
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// إذا كان المعرف غير صالح (ليس رقماً)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "معرف المقال غير صالح."})
	}

	// 3. استدعاء طبقة منطق العمل (الطابق الثاني) لحذف المقال
	if err := h.articleUseCase.DeleteArticle(uint(id)); err != nil { 
		// إذا لم يتم العثور على المقال
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("المقال بالمعرف %d غير موجود.", id)})
		}
		// تسجيل الخطأ للمطورين
		log.Printf("خطأ في حذف المقال: %v", err)
		// إرجاع رسالة خطأ للمستخدم
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "فشل حذف المقال."})
	}

	// 4. إذا نجحت العملية، إرسال رد بأن العملية تمت بنجاح (بدون محتوى)
	return c.Status(fiber.StatusNoContent).SendString("")
}