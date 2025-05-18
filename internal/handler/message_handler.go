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

// StartDispatcher godoc
// @Summary      Start the dispatcher
// @Description  Starts the message dispatcher
// @Tags         Dispatcher
// @Accept       json
// @Produce      plain
// @Success      200 {string} string "Dispatcher started"
// @Router       /start [post]
func (h *MessageHandler) StartDispatcher(c *fiber.Ctx) error {
	h.dispatcher.Start()
	return c.SendString("Dispatcher started")
}

// StopDispatcher godoc
// @Summary      Stop the dispatcher
// @Description  Stops the message dispatcher
// @Tags         Dispatcher
// @Accept       json
// @Produce      plain
// @Success      200 {string} string "Dispatcher stopped"
// @Router       /stop [post]
func (h *MessageHandler) StopDispatcher(c *fiber.Ctx) error {
	h.dispatcher.Stop()
	return c.SendString("Dispatcher stopped")
}

// GetSentMessages godoc
// @Summary      Get sent messages
// @Description  Retrieves the list of sent messages
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Success      200 {array} model.Message "List of sent messages"
// @Failure      500 {string} string "Failed to fetch sent messages"
// @Router       /messages/sent [get]
func (h *MessageHandler) GetSentMessages(c *fiber.Ctx) error {
	ctx := c.Context()

	messages, err := h.repo.GetSentMessages(ctx, 100)
	if err != nil {
		zap.L().Error("Failed to fetch sent messages", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch sent messages")
	}

	return c.JSON(messages)
}
