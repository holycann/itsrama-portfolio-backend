package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/gemini"
)

// MockGeminiHandler is a mock implementation of the GeminiHandler
type MockGeminiHandler struct {
	mock.Mock
}

func (m *MockGeminiHandler) CreateChatSession(c *gin.Context) {
	// Get the request body
	var req gemini.CreateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate required fields
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Return a successful response
	sessionID := uuid.New().String()
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"session_id": sessionID,
		},
		"message": "Chat session created successfully",
	})
}

func (m *MockGeminiHandler) SendMessage(c *gin.Context) {
	// Get the request body
	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate required fields
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"response": "This is a mock response to: " + req.Message,
		},
		"message": "Message processed successfully",
	})
}

func (m *MockGeminiHandler) GenerateEventDescription(c *gin.Context) {
	// Get the request body
	var req struct {
		EventID string `json:"event_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate required fields
	if req.EventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"description": "This is a mock event description for event: " + req.EventID,
		},
		"message": "Event description generated successfully",
	})
}

func TestCreateChatSession(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockHandler := &MockGeminiHandler{}
	router.POST("/api/ai/chat/sessions", mockHandler.CreateChatSession)

	// Test case: successful chat session creation
	t.Run("Successful chat session creation", func(t *testing.T) {
		userID := uuid.New().String()

		createRequest := map[string]interface{}{
			"user_id": userID,
		}

		jsonData, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/chat/sessions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var respBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &respBody)

		// Check if the response data contains the session ID
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data["session_id"])
	})

	// Test case: validation failure
	t.Run("Validation failure - empty user ID", func(t *testing.T) {
		createRequest := map[string]interface{}{
			// UserID intentionally left empty
		}

		jsonData, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/chat/sessions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestSendMessage(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockHandler := &MockGeminiHandler{}
	router.POST("/api/ai/chat/messages", mockHandler.SendMessage)

	// Test case: successful message sending
	t.Run("Successful message sending", func(t *testing.T) {
		message := "Tell me about cultural events in Jakarta"

		sendRequest := map[string]interface{}{
			"message": message,
		}

		jsonData, _ := json.Marshal(sendRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/chat/messages", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var respBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &respBody)

		// Check if the response data contains the AI response
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data["response"])
	})

	// Test case: validation failure
	t.Run("Validation failure - empty message", func(t *testing.T) {
		sendRequest := map[string]interface{}{
			// Message intentionally left empty
		}

		jsonData, _ := json.Marshal(sendRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/chat/messages", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGenerateEventDescription(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockHandler := &MockGeminiHandler{}
	router.POST("/api/ai/generate/event-description", mockHandler.GenerateEventDescription)

	// Test case: successful event description generation
	t.Run("Successful event description generation", func(t *testing.T) {
		eventID := uuid.New().String()

		generateRequest := map[string]interface{}{
			"event_id": eventID,
		}

		jsonData, _ := json.Marshal(generateRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/generate/event-description", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var respBody map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &respBody)

		// Check if the response data contains the generated description
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data["description"])
	})

	// Test case: validation failure
	t.Run("Validation failure - empty event ID", func(t *testing.T) {
		generateRequest := map[string]interface{}{
			// EventID intentionally left empty
		}

		jsonData, _ := json.Marshal(generateRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/ai/generate/event-description", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
