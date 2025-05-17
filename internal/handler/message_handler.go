package handler

import (
	"messaging-app/internal/repository"
	"messaging-app/internal/scheduler"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MessageHandler struct {
	repo       repository.MessageRepository
	dispatcher scheduler.Dispatcher
}

func NewMessageHandler(repo repository.MessageRepository, dispatcher scheduler.Dispatcher) *MessageHandler {
	return &MessageHandler{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (h *MessageHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/start", h.StartDispatcher)
	app.Post("/stop", h.StopDispatcher)
	app.Get("/messages/sent", h.GetSentMessages)
}

func (h *MessageHandler) StartDispatcher(c *fiber.Ctx) error {
	h.dispatcher.Start()
	return c.SendString("Dispatcher started")
}

func (h *MessageHandler) StopDispatcher(c *fiber.Ctx) error {
	h.dispatcher.Stop()
	return c.SendString("Dispatcher stopped")
}

func (h *MessageHandler) GetSentMessages(c *fiber.Ctx) error {
	ctx := c.Context()

	messages, err := h.repo.GetSentMessages(ctx, 100)
	if err != nil {
		zap.L().Error("Failed to fetch sent messages", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch sent messages")
	}

	return c.JSON(messages)
}
