// my-article-app/internal/usecase/article_usecase.go
package usecase

// استيراد المكتبات اللازمة للعمل
import (
	"my-article-app/internal/models"     // استيراد نماذج البيانات (مثل Article)
	"my-article-app/internal/repository" // استيراد طبقة المستودع (Repository)
)
type ArticleUseCase interface {
	CreateArticle(article *models.Article) error
	GetAllArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	UpdateArticle(article *models.Article) error
	DeleteArticle(id uint) error
}
// ArticleUseCase يمثل منطق العمل للمقالات (هذا هو الطابق الثاني الذي يحتوي على قواعد العمل)

type articleUseCase struct {
	articleRepo repository.ArticleRepository // مرجع لمستودع المقالات (الطابق الثالث)
	// يمكنك إضافة AuthorRepository هنا إذا احتجت للتحقق من المؤلف
	// authorRepo  *repository.AuthorRepository
}

// NewArticleUseCase ينشئ مثيلاً جديدًا من ArticleUseCase

func NewArticleUseCase(articleRepo repository.ArticleRepository) ArticleUseCase {
	// إنشاء وإرجاع كائن منطق العمل مع تمرير مرجع لمستودع المقالات
	return &articleUseCase{articleRepo: articleRepo }
}

// CreateArticle يقوم بإنشاء مقال جديد (هنا يمكن إضافة منطق عمل إضافي)
// تستقبل هذه الدالة بيانات المقال من طبقة المعالج (Handler) وتقوم بمعالجتها
func (uc *articleUseCase) CreateArticle(article *models.Article) error {
	//  يمكن إضافة أي منطق عمل قبل الحفظ
	// مثال: التحقق من وجود المؤلف (إذا تم تمرير authorRepo)
	// author, err := uc.authorRepo.FindByID(article.AuthorID)
	// if err != nil || author == nil {
	// 	return errors.New("المؤلف غير موجود")
	// }

	//  يمكن إضافة قواعد عمل أخرى مثل:
	// - التحقق من أن المقال ليس مكرراً
	// - إضافة تاريخ النشر تلقائياً
	// - تنسيق محتوى المقال

	// بعد التحقق من كل قواعد العمل، نطلب من طبقة المستودع (الطابق الثالث) حفظ المقال
	return uc.articleRepo.Create(article)
}

// GetAllArticles يجلب جميع المقالات
// تستدعي هذه الدالة من طبقة المعالج (Handler) عندما يريد المستخدم عرض جميع المقالات
func (uc *articleUseCase) GetAllArticles() ([]models.Article, error) {
	//  يمكن إضافة منطق عمل إضافي مثل:
	// - ترتيب المقالات حسب التاريخ
	// - تصفية المقالات حسب حالة النشر
	// - التحقق من صلاحيات المستخدم

	// طلب جميع المقالات من طبقة المستودع (الطابق الثالث)
	return uc.articleRepo.FindAll()
}

// GetArticleByID يجلب مقالًا واحدًا حسب ID
// تستدعي هذه الدالة من طبقة المعالج (Handler) عندما يريد المستخدم عرض مقال محدد
func (uc *articleUseCase) GetArticleByID(id uint) (*models.Article, error) {
	// هنا يمكن إضافة منطق عمل إضافي مثل:
	// - التحقق من صلاحيات المستخدم لعرض هذا المقال
	// - تسجيل عملية الوصول للمقال في سجل الأحداث
	// - زيادة عداد المشاهدات للمقال

	// طلب المقال المحدد من طبقة المستودع (الطابق الثالث)
	return uc.articleRepo.FindByID(id)
}

// UpdateArticle يقوم بتحديث مقال موجود
// تستدعي هذه الدالة من طبقة المعالج (Handler) عندما يريد المستخدم تحديث مقال
func (uc *articleUseCase) UpdateArticle(article *models.Article) error {
	//  يمكن إضافة منطق عمل قبل التحديث مثل:
	// - التحقق من صلاحيات المستخدم لتحديث المقال
	// - التحقق من أن المقال موجود قبل التحديث
	// - تحديث تاريخ التعديل تلقائياً
	// - الاحتفاظ بنسخة من المقال القديم قبل التحديث

	// طلب تحديث المقال من طبقة المستودع (الطابق الثالث)
	return uc.articleRepo.Update(article)
}

// DeleteArticle يقوم بحذف مقال
// تستدعي هذه الدالة من طبقة المعالج (Handler) عندما يريد المستخدم حذف مقال
func (uc *articleUseCase) DeleteArticle(id uint) error {
	//  يمكن إضافة منطق عمل قبل الحذف مثل:
	// - التحقق من صلاحيات المستخدم لحذف المقال
	// - التحقق من عدم وجود تعليقات أو إعجابات مرتبطة بالمقال
	// - عمل أرشفة للمقال بدلاً من حذفه نهائياً
	// - تسجيل عملية الحذف في سجل الأحداث

	// طلب حذف المقال من طبقة المستودع (الطابق الثالث)
	return uc.articleRepo.Delete(id)
}
