// my-article-app/internal/repository/author_repository.go
package repository

// استيراد المكتبات اللازمة للعمل مع قاعدة البيانات
import (
	"fmt"                        // مكتبة لتنسيق النصوص ورسائل الخطأ
	"my-article-app/internal/models" // استيراد نماذج البيانات (مثل Author)

	"gorm.io/gorm"               // مكتبة GORM للتعامل مع قواعد البيانات
)

type AuthorRepository interface {

	Create(author *models.Author) error
	FindAll() ([]models.Author, error)
	FindByID(id uint) (*models.Author, error)
	Update(author *models.Author) error
	Delete(id uint) error
}

type authorRepository struct {
	db *gorm.DB
}

// NewAuthorRepository ينشئ مثيلاً جديدًا من AuthorRepository

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	// إنشاء وإرجاع كائن مستودع جديد مع تمرير اتصال قاعدة البيانات له
	return &authorRepository{db: db}
}

// Create ينشئ مؤلفًا جديدًا في قاعدة البيانات
// هذه الدالة مسؤولة عن حفظ بيانات مؤلف جديد في قاعدة البيانات
func (r *authorRepository) Create(author *models.Author) error {
	// استخدام GORM لإنشاء سجل جديد في قاعدة البيانات
	// سيتم تعبئة حقل ID تلقائيًا بعد الإنشاء الناجح
	result := r.db.Create(author)

	// التحقق من حدوث أي خطأ أثناء الإنشاء
	if result.Error != nil {
		// إرجاع رسالة خطأ منسقة مع الخطأ الأصلي
		return fmt.Errorf("فشل إنشاء المؤلف: %w", result.Error)
	}
	
	// إرجاع nil في حالة نجاح العملية
	return nil
}

// FindAll يجلب جميع المؤلفين من قاعدة البيانات
// هذه الدالة مسؤولة عن استرجاع جميع سجلات المؤلفين من قاعدة البيانات
func (r *authorRepository) FindAll() ([]models.Author, error) {
	// إنشاء متغير من نوع مصفوفة لتخزين المؤلفين المسترجعين
	var authors []models.Author

	// يمكن استخدام Preload("Articles") هنا إذا أردت جلب المقالات المرتبطة بكل مؤلف
	// مثال: r.db.Preload("Articles").Find(&authors)

	// استخدام GORM لجلب جميع سجلات المؤلفين
	result := r.db.Find(&authors)

	// التحقق من حدوث أي خطأ أثناء الاستعلام
	if result.Error != nil {
		// إرجاع nil ورسالة خطأ في حالة الفشل
		return nil, fmt.Errorf("فشل جلب المؤلفين: %w", result.Error)
	}

	// إرجاع مصفوفة المؤلفين وnil للخطأ في حالة النجاح
	return authors, nil
}

// FindByID يجلب مؤلفًا واحدًا حسب ID
// هذه الدالة تُستخدم لاسترجاع مؤلف معين من قاعدة البيانات باستخدام معرفه الفريد
func (r *authorRepository) FindByID(id uint) (*models.Author, error) {
	// إنشاء متغير من نوع Author لتخزين المؤلف المسترجع
	var author models.Author

	// استخدام Preload("Articles") لجلب المقالات المرتبطة بالمؤلف
	// هذا يعني أننا سنجلب المؤلف مع جميع مقالاته في استعلام واحد
	result := r.db.Preload("Articles").First(&author, id)

	// التحقق من حدوث أي خطأ أثناء الاستعلام
	if result.Error != nil {
		// التحقق مما إذا كان الخطأ هو عدم العثور على السجل
		if result.Error == gorm.ErrRecordNotFound {
			// إرجاع nil, nil للإشارة إلى أنه لم يتم العثور على المؤلف
			return nil, nil // لم يتم العثور على المؤلف
		}
		// إرجاع nil ورسالة خطأ في حالة حدوث خطأ آخر
		return nil, fmt.Errorf("فشل جلب المؤلف بالمعرف %d: %w", id, result.Error)
	}

	// إرجاع مرجع للمؤلف وnil للخطأ في حالة النجاح
	return &author, nil
}

// Update يقوم بتحديث مؤلف موجود في قاعدة البيانات
// هذه الدالة مسؤولة عن تحديث بيانات مؤلف موجود بالفعل في قاعدة البيانات
func (r *authorRepository) Update(author *models.Author) error {
	// استخدام دالة Save من GORM لحفظ التغييرات على المؤلف
	// إذا كان المؤلف موجودًا (له ID) فسيتم تحديثه، وإلا سيتم إنشاؤه
	result := r.db.Save(author)

	// التحقق من حدوث أي خطأ أثناء التحديث
	if result.Error != nil {
		// إرجاع رسالة خطأ منسقة مع الخطأ الأصلي
		return fmt.Errorf("فشل تحديث المؤلف: %w", result.Error)
	}

	// التحقق من أن التحديث أثر على سجل واحد على الأقل
	// إذا كانت الصفوف المتأثرة = 0، فهذا يعني أن المؤلف غير موجود
	if result.RowsAffected == 0 {
		// إرجاع خطأ قياسي من GORM يشير إلى أن السجل غير موجود
		return gorm.ErrRecordNotFound
	}

	// إرجاع nil في حالة نجاح العملية
	return nil
}

// Delete يقوم بحذف مؤلف من قاعدة البيانات باستخدام ID
// هذه الدالة مسؤولة عن حذف مؤلف من قاعدة البيانات باستخدام معرّفه الفريد
func (r *authorRepository) Delete(id uint) error {
	// استخدام GORM لحذف المؤلف بواسطة المعرّف
	// نمرر كائن Author فارغ ومعرّف المؤلف المراد حذفه
	// ملاحظة: اعتمادًا على إعدادات GORM، قد يكون هذا حذفًا فعليًا أو حذفًا منطقيًا (soft delete)
	result := r.db.Delete(&models.Author{}, id)

	// التحقق من حدوث أي خطأ أثناء الحذف
	if result.Error != nil {
		// إرجاع رسالة خطأ منسقة مع الخطأ الأصلي
		return fmt.Errorf("فشل حذف المؤلف: %w", result.Error)
	}

	// التحقق من أن الحذف أثر على سجل واحد على الأقل
	// إذا كانت الصفوف المتأثرة = 0، فهذا يعني أن المؤلف غير موجود
	if result.RowsAffected == 0 {
		// إرجاع خطأ قياسي من GORM يشير إلى أن السجل غير موجود
		return gorm.ErrRecordNotFound
	}

	// إرجاع nil في حالة نجاح العملية
	return nil
}