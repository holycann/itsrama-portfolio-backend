package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/discussion/handlers"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Mock message service
type mockMessageService struct {
	mock.Mock
}

func (m *mockMessageService) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *mockMessageService) GetMessageByID(ctx context.Context, id string) (*models.MessageDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MessageDTO), args.Error(1)
}

func (m *mockMessageService) GetMessagesByThread(ctx context.Context, threadID string) ([]models.MessageDTO, error) {
	args := m.Called(ctx, threadID)
	return args.Get(0).([]models.MessageDTO), args.Error(1)
}

func (m *mockMessageService) GetMessagesByUser(ctx context.Context, userID string) ([]models.MessageDTO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.MessageDTO), args.Error(1)
}

func (m *mockMessageService) GetRecentMessages(ctx context.Context, limit int) ([]models.MessageDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.MessageDTO), args.Error(1)
}

func (m *mockMessageService) ListMessages(ctx context.Context, opts base.ListOptions) ([]models.MessageDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.MessageDTO), args.Error(1)
}

func (m *mockMessageService) SearchMessages(ctx context.Context, query string, opts base.ListOptions) ([]models.MessageDTO, int, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]models.MessageDTO), args.Int(1), args.Error(2)
}

func (m *mockMessageService) CountMessages(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockMessageService) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *mockMessageService) DeleteMessage(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to setup test environment
func setupMessageHandlerTest() (*gin.Engine, *mockMessageService, *handlers.MessageHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockMessageService)
	appLogger := logger.DefaultLogger()
	messageHandler := handlers.NewMessageHandler(mockService, appLogger)

	return router, mockService, messageHandler
}

func TestCreateMessage(t *testing.T) {
	router, mockService, messageHandler := setupMessageHandlerTest()

	// Test case: successful message creation
	t.Run("Successful message creation", func(t *testing.T) {
		router.POST("/messages", messageHandler.CreateMessage)

		messageID := uuid.New()
		threadID := uuid.New()
		senderID := uuid.New()

		now := time.Now()

		newMessage := &models.Message{
			ID:        messageID,
			ThreadID:  threadID,
			SenderID:  senderID,
			Content:   "Test message content",
			Type:      models.DiscussionMessageType,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		mockService.On("CreateMessage", mock.Anything, mock.AnythingOfType("*models.Message")).Return(newMessage, nil)

		createMessage := models.CreateMessage{
			ThreadID: threadID,
			SenderID: senderID,
			Content:  "Test message content",
			Type:     models.DiscussionMessageType,
		}

		jsonData, _ := json.Marshal(createMessage)
		req, _ := http.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetMessagesByThread(t *testing.T) {
	router, mockService, messageHandler := setupMessageHandlerTest()

	// Test case: successful retrieval of messages by thread
	t.Run("Successful retrieval of messages by thread", func(t *testing.T) {
		router.GET("/threads/:id/messages", messageHandler.GetMessagesByThread)

		threadID := uuid.New()
		senderID1 := uuid.New()
		senderID2 := uuid.New()
		messageID1 := uuid.New()
		messageID2 := uuid.New()

		now := time.Now()

		messages := []models.MessageDTO{
			{
				ID:        messageID1,
				ThreadID:  threadID,
				Content:   "First message",
				Type:      models.DiscussionMessageType,
				CreatedAt: &now,
				Sender: &userModels.User{
					ID:    senderID1,
					Email: "user1@example.com",
				},
			},
			{
				ID:        messageID2,
				ThreadID:  threadID,
				Content:   "Second message",
				Type:      models.DiscussionMessageType,
				CreatedAt: &now,
				Sender: &userModels.User{
					ID:    senderID2,
					Email: "user2@example.com",
				},
			},
		}

		mockService.On("GetMessagesByThread", mock.Anything, threadID.String()).Return(messages, nil)

		req, _ := http.NewRequest(http.MethodGet, "/threads/"+threadID.String()+"/messages", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: no messages found for thread
	t.Run("No messages found for thread", func(t *testing.T) {
		router.GET("/threads/:id/messages", messageHandler.GetMessagesByThread)

		threadID := uuid.New()
		var emptyMessages []models.MessageDTO

		mockService.On("GetMessagesByThread", mock.Anything, threadID.String()).Return(emptyMessages, nil)

		req, _ := http.NewRequest(http.MethodGet, "/threads/"+threadID.String()+"/messages", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListMessages(t *testing.T) {
	router, mockService, messageHandler := setupMessageHandlerTest()

	// Test case: successful message listing
	t.Run("Successful message listing", func(t *testing.T) {
		router.GET("/messages", messageHandler.ListMessages)

		threadID := uuid.New()
		senderID := uuid.New()
		messageID1 := uuid.New()
		messageID2 := uuid.New()

		now := time.Now()

		messages := []models.MessageDTO{
			{
				ID:        messageID1,
				ThreadID:  threadID,
				Content:   "First message",
				Type:      models.DiscussionMessageType,
				CreatedAt: &now,
				Sender: &userModels.User{
					ID:    senderID,
					Email: "user@example.com",
				},
			},
			{
				ID:        messageID2,
				ThreadID:  threadID,
				Content:   "Second message",
				Type:      models.DiscussionMessageType,
				CreatedAt: &now,
				Sender: &userModels.User{
					ID:    senderID,
					Email: "user@example.com",
				},
			},
		}

		mockService.On("ListMessages", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(messages, nil)

		req, _ := http.NewRequest(http.MethodGet, "/messages", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateMessage(t *testing.T) {
	router, mockService, messageHandler := setupMessageHandlerTest()

	// Test case: successful message update
	t.Run("Successful message update", func(t *testing.T) {
		router.PUT("/messages/:id", messageHandler.UpdateMessage)

		messageID := uuid.New()
		threadID := uuid.New()
		senderID := uuid.New()

		now := time.Now()

		updatedMessage := &models.Message{
			ID:        messageID,
			ThreadID:  threadID,
			SenderID:  senderID,
			Content:   "Updated message content",
			Type:      models.DiscussionMessageType,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		mockService.On("UpdateMessage", mock.Anything, mock.AnythingOfType("*models.Message")).Return(updatedMessage, nil)

		updatePayload := map[string]interface{}{
			"id":      messageID.String(),
			"content": "Updated message content",
		}

		jsonData, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest(http.MethodPut, "/messages/"+messageID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteMessage(t *testing.T) {
	router, mockService, messageHandler := setupMessageHandlerTest()

	// Test case: successful message deletion
	t.Run("Successful message deletion", func(t *testing.T) {
		router.DELETE("/messages/:id", messageHandler.DeleteMessage)

		messageID := uuid.New()
		mockService.On("DeleteMessage", mock.Anything, messageID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/messages/"+messageID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
