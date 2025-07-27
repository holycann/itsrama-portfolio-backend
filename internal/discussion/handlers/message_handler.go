package handlers

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/holycann/cultour-backend/pkg/repository"
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
// @Success 201 {object} response.APIResponse{data=models.ResponseMessage} "Message created successfully"
// @Failure 400 {object} response.APIResponse "Invalid message creation details"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Validate required fields
	if message.ThreadID == uuid.Nil || message.UserID == uuid.Nil || message.Content == "" {
		// Convert map to JSON string for error details
		details, _ := json.Marshal(map[string]interface{}{
			"thread_id": message.ThreadID == uuid.Nil,
			"user_id":   message.UserID == uuid.Nil,
			"content":   message.Content == "",
		})
		response.BadRequest(c, "Missing required fields", string(details), "")
		return
	}

	if err := h.messageService.CreateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error creating message: %v", err)
		response.InternalServerError(c, "Failed to create message", err.Error(), "")
		return
	}

	response.SuccessCreated(c, message, "Message created successfully")
}

// ListMessages godoc
// @Summary List messages
// @Description Retrieve a list of messages with pagination and filtering
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param limit query int false "Number of messages to retrieve" default(10)
// @Param offset query int false "Number of messages to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.ResponseMessage} "Messages retrieved successfully"
// @Failure 500 {object} response.APIResponse "Failed to list messages"
// @Router /messages [get]
func (h *MessageHandler) ListMessages(c *gin.Context) {
	// Parse pagination parameters with defaults
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Prepare list options
	listOptions := repository.ListOptions{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: repository.SortDescending,
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Optional filtering
	filters := []repository.FilterOption{}
	if threadID := c.Query("thread_id"); threadID != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "thread_id",
			Operator: "=",
			Value:    threadID,
		})
	}
	if userID := c.Query("user_id"); userID != "" {
		filters = append(filters, repository.FilterOption{
			Field:    "user_id",
			Operator: "=",
			Value:    userID,
		})
	}
	listOptions.Filters = filters

	// Retrieve messages
	messages, err := h.messageService.ListMessages(c.Request.Context(), listOptions)
	if err != nil {
		h.logger.Error("Error retrieving messages: %v", err)
		response.InternalServerError(c, "Failed to retrieve messages", err.Error(), "")
		return
	}

	// Count total messages for pagination
	totalMessages, err := h.messageService.CountMessages(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("Error counting messages: %v", err)
		response.InternalServerError(c, "Failed to count messages", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalMessages,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalMessages) / float64(limit))),
		HasNextPage: offset+limit < totalMessages,
	}

	// Respond with messages and pagination
	response.SuccessOK(c, messages, "Messages retrieved successfully", pagination)
}

// SearchMessages godoc
// @Summary Search messages
// @Description Search messages by various criteria
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string false "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search query (content, etc.)"
// @Param limit query int false "Number of results to retrieve" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param sort_by query string false "Field to sort by" default("created_at")
// @Param sort_order query string false "Sort order (asc/desc)" default("desc")
// @Success 200 {object} response.APIResponse{data=[]models.ResponseMessage} "Messages found successfully"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /messages/search [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	// Get search query
	query := c.Query("query")
	if query == "" {
		response.BadRequest(c, "Search query is required", "Empty search query", "")
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Prepare list options for search
	listOptions := repository.ListOptions{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: repository.SortDescending,
		Filters: []repository.FilterOption{
			{
				Field:    "content",
				Operator: "like",
				Value:    query,
			},
		},
	}
	if sortOrder == "asc" {
		listOptions.SortOrder = repository.SortAscending
	}

	// Search messages
	messages, err := h.messageService.SearchMessages(c.Request.Context(), query, listOptions)
	if err != nil {
		h.logger.Error("Error searching messages: %v", err)
		response.InternalServerError(c, "Failed to search messages", err.Error(), "")
		return
	}

	// Count total search results
	totalMessages, err := h.messageService.CountMessages(c.Request.Context(), listOptions.Filters)
	if err != nil {
		h.logger.Error("Error counting search results: %v", err)
		response.InternalServerError(c, "Failed to count search results", err.Error(), "")
		return
	}

	// Create pagination struct
	pagination := &response.Pagination{
		Total:       totalMessages,
		Page:        offset/limit + 1,
		PerPage:     limit,
		TotalPages:  int(math.Ceil(float64(totalMessages) / float64(limit))),
		HasNextPage: offset+limit < totalMessages,
	}

	// Respond with messages and pagination
	response.SuccessOK(c, messages, "Messages found successfully", pagination)
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
// @Success 200 {object} response.APIResponse{data=models.ResponseMessage} "Message updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid message update details"
// @Failure 404 {object} response.APIResponse "Message not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /messages/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	// Get message ID from path parameter
	messageID := c.Param("id")
	if messageID == "" {
		response.BadRequest(c, "Message ID is required", "Missing message ID", "")
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error(), "")
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(messageID)
	if err != nil {
		response.BadRequest(c, "Invalid Message ID", "Invalid UUID format", "")
		return
	}
	message.ID = parsedID

	if err := h.messageService.UpdateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error updating message: %v", err)
		response.InternalServerError(c, "Failed to update message", err.Error(), "")
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
// @Success 200 {object} response.APIResponse "Deleted successfully"
// @Failure 400 {object} response.APIResponse "Invalid message ID"
// @Failure 404 {object} response.APIResponse "Message not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	// Get message ID from path parameter
	messageID := c.Param("id")
	if messageID == "" {
		response.BadRequest(c, "Message ID is required", "Missing message ID", "")
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), messageID); err != nil {
		h.logger.Error("Error deleting message: %v", err)
		response.InternalServerError(c, "Failed to delete message", err.Error(), "")
		return
	}

	response.SuccessOK(c, nil, "Message deleted successfully")
}
