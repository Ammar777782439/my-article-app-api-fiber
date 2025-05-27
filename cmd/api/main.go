// my-article-app/cmd/api/main.go
package main

import (
	"log"
	"my-article-app/internal/database"
	"my-article-app/internal/handlers"
	"my-article-app/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. تهيئة اتصال قاعدة البيانات باستخدام GORM
	db, err := database.InitGORMDB()
	if err != nil {
		log.Fatalf("فشل في تهيئة قاعدة البيانات: %v", err)
	}

	// 2. تهيئة الـ Repositories (المستودعات)
	articleRepo := repository.NewArticleRepository(db)
	// تهيئة مستودع المؤلفين
	authorRepo := repository.NewAuthorRepository(db) // تحتاج لإنشاء هذا الريبو أولاً

	// 3. تهيئة الـ Handlers (المعالجات)
	// الآن، NewArticleHandler قد يحتاج إلى AuthorRepository إذا كنت تقوم بالتحقق من وجود المؤلف
	articleHandler := handlers.NewArticleHandler(articleRepo /*, authorRepo*/) // تمرير AuthorRepo إذا لزم الأمر
	// تهيئة معالج المؤلفين
	authorHandler := handlers.NewAuthorHandler(authorRepo) // تحتاج لإنشاء هذا الـ Handler أولاً

	app := fiber.New()

	// 4. تعريف مسارات Fiber (Routes)
	api := app.Group("/api/v1") // مجموعة (Group) للـ API

	// مسارات المقالات
	articlesGroup := api.Group("/articles") // مجموعة خاصة بالمقالات
	articlesGroup.Post("/", articleHandler.CreateArticle)
	articlesGroup.Get("/", articleHandler.GetAllArticles)
	articlesGroup.Get("/:id", articleHandler.GetArticleByID)
	articlesGroup.Put("/:id", articleHandler.UpdateArticle)
	articlesGroup.Delete("/:id", articleHandler.DeleteArticle)

	// مسارات المؤلفين (الجديدة)
	authorsGroup := api.Group("/authors") // مجموعة خاصة بالمؤلفين
	authorsGroup.Post("/", authorHandler.CreateAuthor)
	authorsGroup.Get("/", authorHandler.GetAllAuthors)
	authorsGroup.Get("/:id", authorHandler.GetAuthorByID)
	authorsGroup.Put("/:id", authorHandler.UpdateAuthor)
	authorsGroup.Delete("/:id", authorHandler.DeleteAuthor)

	// مسار للتحقق من حالة الخادم (Health Check)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Application is healthy!")
	})

	// مسار الـ 404 (يجب أن يبقى في النهاية)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "عذراً، المسار غير موجود (404)."})
	})

	// تشغيل الخادم على المنفذ 3000
	log.Fatal(app.Listen(":3000"))
}