package main

import (
	"database/sql"
	"log"
	"path/filepath"

	"messaging-app/internal/cache"
	"messaging-app/internal/handler"
	"messaging-app/internal/logger"
	"messaging-app/internal/repository"
	"messaging-app/internal/scheduler"
	"messaging-app/internal/service"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Logger setup
	logger := logger.CreateLogger()
	zap.ReplaceGlobals(logger)

	// PostgreSQL bağlantısı
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/messaging?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Redis bağlantısı
	redisCache := cache.NewRedisCache("localhost:6379")

	// Katmanların kurulması
	repo := repository.NewMessageRepository(db)
	svc := service.NewMessageService(repo, redisCache)
	dispatcher := scheduler.NewDispatcher(svc)
	h := handler.NewMessageHandler(repo, dispatcher)

	// Fiber v2 app başlat
	app := fiber.New()
	h.RegisterRoutes(app)

	// Swagger UI static dosyalarını sun
	swaggerDir, _ := filepath.Abs("../static/swagger")
	app.Static("/swagger", swaggerDir)

	// Swagger index.html
	app.Get("/swagger", func(c *fiber.Ctx) error {
		indexPath, err := filepath.Abs("../static/swagger/index.html")
		if err != nil {
			return err
		}
		return c.SendFile(indexPath, true)
	})

	// swagger.yaml dosyasını sun
	app.Get("/swagger.yaml", func(c *fiber.Ctx) error {
		yamlPath, err := filepath.Abs("../docs/swagger.yaml")
		if err != nil {
			return err
		}
		return c.SendFile(yamlPath, true)
	})

	log.Println("Fiber server is running on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start Fiber app:", err)
	}
}
