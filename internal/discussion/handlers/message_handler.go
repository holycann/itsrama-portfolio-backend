package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
)

// MessageHandler handles HTTP requests related to messages
type MessageHandler struct {
	messageService services.MessageService
	logger         *logger.Logger
}

// NewMessageHandler creates a new instance of message handler
func NewMessageHandler(messageService services.MessageService, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		logger:         logger,
	}
}

// CreateMessage godoc
// @Summary Create a new message
// @Description Add a new message to the system
// @Tags Messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param message body models.Message true "Message Information"
// @Success 201 {object} response.Response{data=models.ResponseMessage} "Message created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid message creation details"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.messageService.CreateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error creating message: %v", err)
		response.InternalServerError(c, "Failed to create message", err.Error())
		return
	}

	response.SuccessCreated(c, message, "Message created successfully")
}

// ListMessages godoc
// @Summary List messages
// @Description Retrieve a list of messages with pagination
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of messages to retrieve" default(10)
// @Param offset query int false "Number of messages to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseMessage} "Messages retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to list messages"
// @Router /messages [get]
func (h *MessageHandler) ListMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parsing limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// Get list of messages
	messages, err := h.messageService.GetMessages(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving messages: %v", err)
		response.InternalServerError(c, "Failed to retrieve messages", err.Error())
		return
	}

	// Count total messages for pagination
	total, err := h.messageService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting messages: %v", err)
		response.InternalServerError(c, "Failed to count messages", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, messages, total, offset/limit+1, limit)
}

// SearchMessages godoc
// @Summary Search messages
// @Description Search messages by various criteria
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id query string false "Message ID"
// @Param threadID query string false "Thread ID"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseMessage} "Messages found successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /messages/search [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	// Get query parameters
	id := c.Query("id")
	threadID := c.Query("threadID")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parsing limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// If ID is provided, search by ID
	if id != "" {
		message, err := h.messageService.GetMessageByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding message by ID: %v", err)
			response.NotFound(c, "Message not found", err.Error())
			return
		}
		response.SuccessOK(c, message, "Message found")
		return
	}

	// If threadID is provided, search by threadID
	if threadID != "" {
		messages, err := h.messageService.GetMessagesByThreadID(c.Request.Context(), threadID, limit, offset)
		if err != nil {
			h.logger.Error("Error finding messages by thread ID: %v", err)
			response.NotFound(c, "Messages not found", err.Error())
			return
		}

		// Count total messages for pagination
		total, err := h.messageService.CountByThreadID(c.Request.Context(), threadID)
		if err != nil {
			h.logger.Error("Error counting messages: %v", err)
			response.InternalServerError(c, "Failed to count messages", err.Error())
			return
		}

		response.WithPagination(c, messages, total, offset/limit+1, limit)
		return
	}

	// If no specific parameters are provided, return a list of messages
	messages, err := h.messageService.GetMessages(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving messages: %v", err)
		response.InternalServerError(c, "Failed to retrieve messages", err.Error())
		return
	}

	// Count total messages for pagination
	total, err := h.messageService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting messages: %v", err)
		response.InternalServerError(c, "Failed to count messages", err.Error())
		return
	}

	// Use WithPagination to add pagination metadata
	response.WithPagination(c, messages, total, offset/limit+1, limit)
}

// UpdateMessage godoc
// @Summary Update a message
// @Description Update an existing message's details
// @Tags Messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Message ID"
// @Param message body models.Message true "Message Update Details"
// @Success 200 {object} response.Response{data=models.ResponseMessage} "Message updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid message update details"
// @Failure 404 {object} response.ErrorResponse "Message not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /messages/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.messageService.UpdateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error updating message: %v", err)
		response.InternalServerError(c, "Failed to update message", err.Error())
		return
	}

	response.SuccessOK(c, message, "Message updated successfully")
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Remove a message from the system by its unique identifier
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Message ID"
// @Success 200 {object} response.Response "Message deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid message ID"
// @Failure 404 {object} response.ErrorResponse "Message not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Message ID is required", nil)
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting message: %v", err)
		response.InternalServerError(c, "Failed to delete message", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Message deleted successfully")
}
