// my-article-app/internal/repository/article_repository.go
package repository

import (
	"fmt"
	"my-article-app/internal/models" // استيراد بنية Article

	"gorm.io/gorm" // استيراد GORM
)

// ArticleRepository يمثل المستودع الذي يتعامل مع عمليات CRUD للمقالات
type ArticleRepository struct {
	db *gorm.DB // لاحظ أن النوع هنا أصبح *gorm.DB بدلاً من *sql.DB
}

// NewArticleRepository ينشئ مثيلاً جديدًا من ArticleRepository
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create ينشئ مقالًا جديدًا في قاعدة البيانات
func (r *ArticleRepository) Create(article *models.Article) error {
	// GORM: db.Create(&article) سيقوم بإنشاء سجل جديد في جدول articles
	// وسيتم ملء حقل ID تلقائياً بواسطة GORM بعد الإنشاء.
	result := r.db.Create(article)
	if result.Error != nil {
		return fmt.Errorf("فشل إنشاء المقال: %w", result.Error)
	}
	return nil
}

// FindAll يجلب جميع المقالات من قاعدة البيانات
func (r *ArticleRepository) FindAll() ([]models.Article, error) {
	var articles []models.Article
	// GORM: db.Find(&articles) سيقوم بجلب جميع السجلات من جدول articles
	result := r.db.Find(&articles)
	if result.Error != nil {
		return nil, fmt.Errorf("فشل جلب المقالات: %w", result.Error)
	}
	return articles, nil
}

// FindByID يجلب مقالًا واحدًا حسب ID
func (r *ArticleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	// GORM: db.First(&article, id) سيقوم بجلب أول سجل يطابق ID المعطى
	result := r.db.First(&article, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound { // GORM له خطأ خاص لعدم وجود السجل
			return nil, nil // لم يتم العثور على المقال
		}
		return nil, fmt.Errorf("فشل جلب المقال بالمعرف %d: %w", id, result.Error)
	}
	return &article, nil
}

// Update يقوم بتحديث مقال موجود في قاعدة البيانات
func (r *ArticleRepository) Update(article *models.Article) error {
	// GORM: db.Save(&article) سيقوم بتحديث السجل إذا كان له ID موجود
	// وإلا فسيقوم بإنشاء سجل جديد (Upsert)
	result := r.db.Save(article) // Save يعمل كـ Update أو Create
	if result.Error != nil {
		return fmt.Errorf("فشل تحديث المقال: %w", result.Error)
	}
	if result.RowsAffected == 0 { // للتحقق مما إذا كان السجل موجوداً وتم تحديثه فعلاً
		// يمكننا أن نفترض هنا أنه إذا لم يتأثر أي صف، فالمقال غير موجود.
		// أو يمكننا جلب المقال أولاً للتأكد.
		// من أجل البساطة هنا، نعتبر أنه لم يتم العثور عليه إذا لم يتأثر أي صف.
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete يقوم بحذف مقال من قاعدة البيانات باستخدام ID
func (r *ArticleRepository) Delete(id uint) error {
	// GORM: db.Delete(&models.Article{}, id) سيقوم بحذف سجل من جدول articles
	// يطابق ID المعطى.
	result := r.db.Delete(&models.Article{}, id)
	if result.Error != nil {
		return fmt.Errorf("فشل حذف المقال: %w", result.Error)
	}
	if result.RowsAffected == 0 { // إذا لم يتأثر أي صف، فالمقال غير موجود
		return gorm.ErrRecordNotFound
	}
	return nil
}