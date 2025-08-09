package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/holycann/cultour-backend/pkg/response"
)

// MessageHandler handles HTTP requests related to messages
type MessageHandler struct {
	base.BaseHandler
	messageService services.MessageService
}

// NewMessageHandler creates a new instance of message handler
func NewMessageHandler(
	messageService services.MessageService,
	logger *logger.Logger,
) *MessageHandler {
	return &MessageHandler{
		BaseHandler:    *base.NewBaseHandler(logger),
		messageService: messageService,
	}
}

// CreateMessage godoc
// @Summary Create a new message
// @Description Allows authenticated users to send a message in a specific discussion thread
// @Description Supports creating different types of messages (discussion, AI-generated)
// @Tags Messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param message body models.CreateMessage true "Message Creation Details"
// @Success 201 {object} response.APIResponse{data=models.MessageDTO} "Message successfully created"
// @Failure 400 {object} response.APIResponse "Invalid message creation payload or validation error"
// @Failure 401 {object} response.APIResponse "Authentication required - missing or invalid token"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient thread access privileges"
// @Failure 500 {object} response.APIResponse "Internal server error during message creation"
// @Router /messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.HandleError(c, errors.New(errors.ErrValidation, "Invalid request payload", err))
		return
	}

	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid user ID format",
			err,
		))
		return
	}

	message.SenderID = parsedUserID
	message.Type = "discussion"

	// Validate required fields
	if message.ThreadID == uuid.Nil || message.SenderID == uuid.Nil || message.Content == "" {
		details, _ := json.Marshal(map[string]interface{}{
			"thread_id": message.ThreadID == uuid.Nil,
			"sender_id": message.SenderID == uuid.Nil,
			"content":   message.Content == "",
		})
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Missing required fields",
			nil,
			errors.WithContext("details", string(details)),
		))
		return
	}

	createdMessage, err := h.messageService.CreateMessage(c.Request.Context(), &message)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessCreated(c, createdMessage, "Message created successfully")
}

// ListMessages godoc
// @Summary Retrieve messages list
// @Description Fetches a paginated list of messages with optional filtering and sorting
// @Description Supports advanced querying with flexible pagination and filtering options
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of messages per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort messages by" default("created_at)" Enum(created_at,sender_id)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Param thread_id query string false "Filter messages by specific thread"
// @Param sender_id query string false "Filter messages by specific sender"
// @Success 200 {object} response.APIResponse{data=[]models.MessageDTO} "Successfully retrieved messages list"
// @Success 204 {object} response.APIResponse "No messages found"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during messages retrieval"
// @Router /messages [get]
func (h *MessageHandler) ListMessages(c *gin.Context) {
	// Parse pagination parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate pagination parameters
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	// Prepare list options
	listOptions := base.ListOptions{
		Page:      page,
		PerPage:   perPage,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// Optional filtering
	filters := []base.FilterOption{}
	if threadID := c.Query("thread_id"); threadID != "" {
		filters = append(filters, base.FilterOption{
			Field:    "thread_id",
			Operator: base.OperatorEqual,
			Value:    threadID,
		})
	}
	if senderID := c.Query("sender_id"); senderID != "" {
		filters = append(filters, base.FilterOption{
			Field:    "sender_id",
			Operator: base.OperatorEqual,
			Value:    senderID,
		})
	}
	listOptions.Filters = filters

	// Retrieve messages
	messages, err := h.messageService.ListMessages(c.Request.Context(), listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Create pagination struct
	data, pagination := base.PaginateResults(messages, listOptions.PerPage, listOptions.Page)

	// Respond with messages and pagination
	response.SuccessOK(c, data, "Messages retrieved successfully", response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}

// SearchMessages godoc
// @Summary Search messages
// @Description Performs a full-text search across message content with advanced filtering
// @Description Allows finding messages by keywords and other attributes
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param query query string true "Search term for finding messages" minlength(2)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param per_page query int false "Number of search results per page" default(10) minimum(1) maximum(100)
// @Param sort_by query string false "Field to sort search results" default("relevance)" Enum(relevance,created_at)
// @Param sort_order query string false "Sort direction" default("desc)" Enum(asc,desc)
// @Success 200 {object} response.APIResponse{data=[]models.MessageDTO} "Successfully completed message search"
// @Success 204 {object} response.APIResponse "No messages match the search query"
// @Failure 400 {object} response.APIResponse "Invalid search parameters"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 500 {object} response.APIResponse "Internal server error during message search"
// @Router /messages/search [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	// Get search query
	query := c.Query("query")
	if query == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Search query is required",
			nil,
		))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate pagination parameters
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	// Prepare list options for search
	listOptions := base.ListOptions{
		Page:      page,
		PerPage:   perPage,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters: []base.FilterOption{
			{
				Field:    "content",
				Operator: base.OperatorLike,
				Value:    query,
			},
		},
	}

	// Search messages
	messages, _, err := h.messageService.SearchMessages(c.Request.Context(), query, listOptions)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Create pagination struct
	data, pagination := base.PaginateResults(messages, listOptions.PerPage, listOptions.Page)

	// Respond with messages and pagination
	response.SuccessOK(c, data, "Messages found successfully", response.WithPagination(pagination.Total, pagination.Page, pagination.PerPage))
}

// UpdateMessage godoc
// @Summary Update an existing message
// @Description Allows message sender to modify their own message content
// @Description Supports partial updates with message type preservation
// @Tags Messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Message Identifier" format(uuid)
// @Param message body models.CreateMessage true "Message Update Payload"
// @Success 200 {object} response.APIResponse{data=models.MessageDTO} "Message successfully updated"
// @Failure 400 {object} response.APIResponse "Invalid message update payload or ID"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - can only update own messages"
// @Failure 404 {object} response.APIResponse "Message not found"
// @Failure 500 {object} response.APIResponse "Internal server error during message update"
// @Router /messages/{thread_id}/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	// Get message ID from path parameter
	messageID := c.Param("id")
	if messageID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Message ID is required",
			nil,
		))
		return
	}

	// Extract user ID from context
	userID, _, _, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrAuthentication,
			"Failed to retrieve user context",
			err,
		))
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.HandleError(c, errors.New(errors.ErrValidation, "Invalid request payload", err))
		return
	}

	// Set the ID from path parameter
	parsedID, err := uuid.Parse(messageID)
	if err != nil {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Invalid Message ID",
			err,
		))
		return
	}
	message.ID = parsedID
	message.Type = "discussion"
	message.SenderID, err = h.ValidateUUID(userID, "User ID")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	updatedMessage, err := h.messageService.UpdateMessage(c.Request.Context(), &message)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, updatedMessage, "Message updated successfully")
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Allows message sender or thread administrator to remove a specific message
// @Description Permanently deletes the message from the discussion thread
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param id path string true "Unique Message Identifier" format(uuid)
// @Success 200 {object} response.APIResponse "Message successfully deleted"
// @Failure 400 {object} response.APIResponse "Invalid message ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient deletion privileges"
// @Failure 404 {object} response.APIResponse "Message not found"
// @Failure 500 {object} response.APIResponse "Internal server error during message deletion"
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	// Get message ID from path parameter
	messageID := c.Param("id")
	if messageID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Message ID is required",
			nil,
		))
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), messageID); err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, nil, "Message deleted successfully")
}

// GetMessagesByThread godoc
// @Summary Retrieve messages for a specific thread
// @Description Fetches all messages associated with a particular discussion thread
// @Description Returns messages in chronological order, supporting thread context
// @Tags Messages
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token (without 'Bearer ' prefix)"
// @Param thread_id path string true "Unique Thread Identifier" format(uuid)
// @Success 200 {object} response.APIResponse{data=[]models.MessageDTO} "Successfully retrieved thread messages"
// @Success 204 {object} response.APIResponse "No messages found in the thread"
// @Failure 400 {object} response.APIResponse "Invalid thread ID format"
// @Failure 401 {object} response.APIResponse "Authentication required"
// @Failure 403 {object} response.APIResponse "Forbidden - insufficient thread access privileges"
// @Failure 404 {object} response.APIResponse "Thread not found"
// @Failure 500 {object} response.APIResponse "Internal server error during message retrieval"
// @Router /messages/thread/{thread_id} [get]
func (h *MessageHandler) GetMessagesByThread(c *gin.Context) {
	// Get thread ID from path parameter
	threadID := c.Param("thread_id")
	if threadID == "" {
		h.HandleError(c, errors.New(
			errors.ErrValidation,
			"Thread ID is required",
			nil,
		))
		return
	}

	// Retrieve messages for the specified thread
	messages, err := h.messageService.GetMessagesByThread(c.Request.Context(), threadID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	response.SuccessOK(c, messages, "Messages retrieved successfully")
}
