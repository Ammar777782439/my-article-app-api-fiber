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

	// 2. تهيئة الـ Repository (المستودع)
	articleRepo := repository.NewArticleRepository(db)

	// 3. تهيئة الـ Handler (المعالج)
	articleHandler := handlers.NewArticleHandler(articleRepo)

	app := fiber.New()

	// 4. تعريف مسارات Fiber (Routes)
	// مجموعة (Group) للمقالات لتنظيم المسارات
	api := app.Group("/api/v1") // يمكننا تجميع المسارات تحت بادئة (prefix) معينة

	articlesGroup := api.Group("/articles") // مجموعة خاصة بالمقالات
	articlesGroup.Post("/", articleHandler.CreateArticle)
	articlesGroup.Get("/", articleHandler.GetAllArticles)
	articlesGroup.Get("/:id", articleHandler.GetArticleByID)
	articlesGroup.Put("/:id", articleHandler.UpdateArticle)
	articlesGroup.Delete("/:id", articleHandler.DeleteArticle)

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