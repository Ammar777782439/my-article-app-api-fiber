// my-article-app/internal/repository/article_repository.go
package repository

// استيراد المكتبات اللازمة للعمل
import (
	// مكتبة للتعامل مع الأخطاء
	"fmt"                            // مكتبة للتعامل مع النصوص
	"my-article-app/internal/models" // استيراد نماذج البيانات (مثل Article)

	"gorm.io/gorm" // مكتبة GORM للتعامل مع قواعد البيانات
)

type ArticleRepository interface {

	Create(article *models.Article) error
	FindAll() ([]models.Article, error)
	FindByID(id uint) (*models.Article, error)
	Update(article *models.Article) error
	Delete(id uint) error
}

type articleRepository struct {
	db *gorm.DB // مرجع لاتصال قاعدة البيانات GORM
}

// NewArticleRepository ينشئ مثيلاً جديدًا من ArticleRepository

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	
	return &articleRepository{db: db}
}

// Create يقوم بإنشاء مقال جديد في قاعدة البيانات
// تُستدعى هذه الدالة من طبقة منطق العمل (UseCase) عندما يُطلب إنشاء مقال جديد
func (r *articleRepository) Create(article *models.Article) error {
	// استخدام GORM لإدخال بيانات المقال في قاعدة البيانات
	// GORM: db.Create(&article) سيقوم بإنشاء سجل جديد في جدول articles
	// وسيتم ملء حقل ID تلقائياً بواسطة GORM بعد الإنشاء.
	// سيقوم GORM أيضاً بحفظ AuthorID إذا تم توفيره في بنية Article
	result := r.db.Create(article)
	if result.Error != nil {
		// إرجاع الخطأ مع رسالة توضيحية
		return fmt.Errorf("فشل إنشاء المقال: %w", result.Error)
	}
	return nil
}

// FindAll يجلب جميع المقالات من قاعدة البيانات
// تُستدعى هذه الدالة من طبقة منطق العمل (UseCase) عندما يُطلب عرض جميع المقالات
func (r *articleRepository) FindAll() ([]models.Article, error) {
	// تعريف متغير لتخزين المقالات المسترجعة
	var articles []models.Article
	// استخدام GORM لاسترجاع جميع المقالات من قاعدة البيانات
	// استخدام Preload("Author") لجلب بيانات المؤلف المرتبطة مع كل مقال
	result := r.db.Preload("Author").Find(&articles)
	if result.Error != nil {
		// إرجاع الخطأ مع رسالة توضيحية
		return nil, fmt.Errorf("فشل جلب المقالات: %w", result.Error)
	}
	 return nil, result.Error
}

// FindByID يجلب مقالًا واحدًا حسب ID
// تُستدعى هذه الدالة من طبقة منطق العمل (UseCase) عندما يُطلب عرض مقال محدد
func (r *articleRepository) FindByID(id uint) (*models.Article, error) {
	// تعريف متغير لتخزين المقال المسترجع
	var article models.Article
	// استخدام GORM للبحث عن المقال بواسطة الـ ID
	// استخدام Preload("Author") لجلب بيانات المؤلف المرتبطة مع المقال
	result := r.db.Preload("Author").First(&article, id)
	if result.Error != nil {
		// إذا كان الخطأ هو عدم وجود المقال
		if result.Error == gorm.ErrRecordNotFound {
			// إرجاع رسالة خطأ مخصصة
			 return nil, result.Error
		}
		// إرجاع الخطأ مع رسالة توضيحية
		return nil, fmt.Errorf("فشل جلب المقال بالمعرف %d: %w", id, result.Error)
	}
	return &article, nil
}

// Update يقوم بتحديث مقال موجود في قاعدة البيانات
// تُستدعى هذه الدالة من طبقة منطق العمل (UseCase) عندما يُطلب تحديث مقال
func (r *articleRepository) Update(article *models.Article) error {
	// استخدام GORM لتحديث السجل إذا كان له ID موجود
	// وإلا فسيقوم بإنشاء سجل جديد (Upsert)
	result := r.db.Save(article) // Save يعمل كـ Update أو Create
	if result.Error != nil {
		// إرجاع الخطأ مع رسالة توضيحية
		return fmt.Errorf("فشل تحديث المقال: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		// للتحقق مما إذا كان السجل موجوداً وتم تحديثه فعلاً
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete يقوم بحذف مقال من قاعدة البيانات باستخدام ID
// تُستدعى هذه الدالة من طبقة منطق العمل (UseCase) عندما يُطلب حذف مقال
func (r *articleRepository) Delete(id uint) error {
	// استخدام GORM لحذف سجل من جدول articles
	// يطابق ID المعطى.
	result := r.db.Delete(&models.Article{}, id)
	if result.Error != nil {
		// إرجاع الخطأ مع رسالة توضيحية
		return fmt.Errorf("فشل حذف المقال: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		// إذا لم يتأثر أي صف، فالمقال غير موجود
		return gorm.ErrRecordNotFound
	}
	return nil
}
