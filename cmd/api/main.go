// my-article-app/cmd/api/main.go
package main

import (
	"log"
	"my-article-app/internal/database"
	"my-article-app/internal/handlers"
	"my-article-app/internal/repository"
	"my-article-app/internal/usecase" // استيراد UseCase

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
	authorRepo := repository.NewAuthorRepository(db)

	// 3. تهيئة الـ Use Cases (حالات الاستخدام) - جديد
	articleUseCase := usecase.NewArticleUseCase(articleRepo)
	authorUseCase := usecase.NewAuthorUseCase(authorRepo)

	// 4. تهيئة الـ Handlers (المعالجات) - استخدام Use Cases
	articleHandler := handlers.NewArticleHandler(articleUseCase) // <-- تغيير هنا
	authorHandler := handlers.NewAuthorHandler(authorUseCase)     // <-- تغيير هنا

	app := fiber.New()

	// 5. تعريف مسارات Fiber (Routes) - لا تغيير هنا
	api := app.Group("/api/v1")

	articlesGroup := api.Group("/articles")
	articlesGroup.Post("/", articleHandler.CreateArticle)
	articlesGroup.Get("/", articleHandler.GetAllArticles)
	articlesGroup.Get("/:id", articleHandler.GetArticleByID)
	articlesGroup.Put("/:id", articleHandler.UpdateArticle)
	articlesGroup.Delete("/:id", articleHandler.DeleteArticle)

	authorsGroup := api.Group("/authors")
	authorsGroup.Post("/", authorHandler.CreateAuthor)
	authorsGroup.Get("/", authorHandler.GetAllAuthors)
	authorsGroup.Get("/:id", authorHandler.GetAuthorByID)
	authorsGroup.Put("/:id", authorHandler.UpdateAuthor)
	authorsGroup.Delete("/:id", authorHandler.DeleteAuthor)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Application is healthy!")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "عذراً، المسار غير موجود (404)."})
	})

	log.Fatal(app.Listen(":3000"))
}