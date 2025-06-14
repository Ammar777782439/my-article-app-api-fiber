// my-article-app/internal/usecase/author_usecase.go
package usecase

// استيراد المكتبات اللازمة للعمل
import (
	"my-article-app/internal/models"
	"my-article-app/internal/repository"
)
type AuthorUseCase interface {
	CreateAuthor(author *models.Author) error
	GetAllAuthors() ([]models.Author, error)
	GetAuthorByID(id uint) (*models.Author, error)
	UpdateAuthor(author *models.Author) error
	DeleteAuthor(id uint) error
}
// AuthorUseCase يمثل منطق العمل للمؤلفين

type authorUseCase struct {
	authorRepo repository.AuthorRepository // مرجع لمستودع المؤلفين (الطبقة الثالثة)
}

// NewAuthorUseCase ينشئ مثيلاً جديدًا من AuthorUseCase
// هذه الدالة تُستخدم عند بدء التطبيق لربط طبقة منطق العمل (UseCase) بطبقة المستودع (Repository)
func NewAuthorUseCase(authorRepo repository.AuthorRepository) AuthorUseCase {
	// إنشاء وإرجاع كائن منطق العمل مع تمرير مرجع لمستودع المؤلفين
	return &authorUseCase{authorRepo: authorRepo}
}

// CreateAuthor يقوم بإنشاء مؤلف جديد
// هذه الدالة مسؤولة عن إنشاء سجل جديد للمؤلف في قاعدة البيانات
func (uc *authorUseCase) CreateAuthor(author *models.Author) error {
	// يمكن إضافة منطق عمل إضافي هنا قبل حفظ المؤلف
	// مثال: التحقق من أن البريد الإلكتروني غير مستخدم بالفعل
	// مثال: التحقق من صحة البيانات المدخلة

	// استدعاء طبقة المستودع لحفظ بيانات المؤلف في قاعدة البيانات
	return uc.authorRepo.Create(author)
}

// GetAllAuthors يجلب جميع المؤلفين
// هذه الدالة مسؤولة عن استرجاع قائمة بجميع المؤلفين من قاعدة البيانات
func (uc *authorUseCase) GetAllAuthors() ([]models.Author, error) {
	// يمكن إضافة منطق عمل إضافي هنا
	// مثال: تطبيق تصفية (filtering) على النتائج
	// مثال: تطبيق ترتيب معين (sorting)

	// استدعاء طبقة المستودع لجلب جميع المؤلفين من قاعدة البيانات
	return uc.authorRepo.FindAll()
}

// GetAuthorByID يجلب مؤلفًا واحدًا حسب ID
// هذه الدالة مسؤولة عن استرجاع مؤلف محدد من قاعدة البيانات باستخدام معرّفه الفريد
func (uc *authorUseCase) GetAuthorByID(id uint) (*models.Author, error) {
	// يمكن إضافة منطق عمل إضافي هنا
	// مثال: التحقق من صلاحية المعرّف قبل البحث

	// استدعاء طبقة المستودع للبحث عن المؤلف بواسطة المعرّف
	// سيتم إرجاع خطأ إذا لم يتم العثور على المؤلف
	return uc.authorRepo.FindByID(id)
}

// UpdateAuthor يقوم بتحديث مؤلف موجود
// هذه الدالة مسؤولة عن تحديث بيانات مؤلف موجود بالفعل في قاعدة البيانات
func (uc *authorUseCase) UpdateAuthor(author *models.Author) error {
	// يمكن إضافة منطق عمل إضافي هنا قبل تحديث المؤلف
	// مثال: التحقق من صحة البيانات المحدثة
	// مثال: التحقق من وجود المؤلف قبل محاولة التحديث

	// استدعاء طبقة المستودع لتحديث بيانات المؤلف في قاعدة البيانات
	return uc.authorRepo.Update(author)
}

// DeleteAuthor يقوم بحذف مؤلف
// هذه الدالة مسؤولة عن حذف مؤلف من قاعدة البيانات باستخدام معرّفه الفريد
func (uc *authorUseCase) DeleteAuthor(id uint) error {
	// يمكن إضافة منطق عمل إضافي هنا قبل حذف المؤلف
	// مثال: التحقق مما إذا كان للمؤلف مقالات مرتبطة به قبل الحذف
	// مثال: تطبيق قواعد العمل مثل حذف المقالات المرتبطة أو منع الحذف إذا كان هناك مقالات

	// استدعاء طبقة المستودع لحذف المؤلف من قاعدة البيانات
	return uc.authorRepo.Delete(id)
}