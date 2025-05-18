package main

import (
	"database/sql"
	"log"
	"os"

	"messaging-app/internal/cache"
	"messaging-app/internal/handler"
	"messaging-app/internal/logger"
	"messaging-app/internal/repository"
	"messaging-app/internal/scheduler"
	"messaging-app/internal/service"

	_ "messaging-app/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger := logger.CreateLogger()
	zap.ReplaceGlobals(logger)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	redisCache := cache.NewRedisCache(os.Getenv("REDIS_ADDR"))

	repo := repository.NewMessageRepository(db)
	svc := service.NewMessageService(repo, redisCache)
	dispatcher := scheduler.NewDispatcher(svc)
	h := handler.NewMessageHandler(repo, dispatcher)

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	h.RegisterRoutes(app)
	log.Println("Fiber server is running on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start Fiber app:", err)
	}
}
