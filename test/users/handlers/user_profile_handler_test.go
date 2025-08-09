package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/handlers"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	errLib "github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserProfileService is a mock implementation of UserProfileService
type MockUserProfileService struct {
	mock.Mock
}

func (m *MockUserProfileService) CreateProfile(ctx context.Context, userProfile *models.UserProfileCreate) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, userProfile)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) UpdateProfile(ctx context.Context, userProfile *models.UserProfileUpdate) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, userProfile)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) UpdateProfileAvatar(ctx context.Context, payload *models.UserProfileAvatarUpdate) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) UpdateProfileIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) DeleteProfile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserProfileService) GetProfileByID(ctx context.Context, id string) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) GetProfileByFullname(ctx context.Context, fullname string) ([]models.UserProfileDTO, error) {
	args := m.Called(ctx, fullname)
	return args.Get(0).([]models.UserProfileDTO), args.Error(1)
}

func (m *MockUserProfileService) ListProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserProfileDTO), args.Int(1), args.Error(2)
}

func (m *MockUserProfileService) SearchProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserProfileDTO), args.Int(1), args.Error(2)
}

func (m *MockUserProfileService) CountProfiles(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func setupUserProfileTestRouter(userProfileHandler *handlers.UserProfileHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware to set user context
	router.Use(func(c *gin.Context) {
		userID := uuid.New().String()
		c.Set("user_id", userID)
		c.Set("email", "test@example.com")
		c.Set("role", "authenticated")
		c.Set("token", "test_token")
		c.Next()
	})

	router.POST("/users/profiles", userProfileHandler.CreateUserProfile)
	router.GET("/users/profiles", userProfileHandler.ListUsersProfile)
	router.GET("/users/profiles/search", userProfileHandler.SearchUserProfile)
	router.GET("/users/profiles/:id", userProfileHandler.GetUserProfileById)
	router.PUT("/users/profiles/:id", userProfileHandler.UpdateUserProfile)
	router.PUT("/users/profiles/:id/avatar", userProfileHandler.UpdateUserAvatar)
	router.DELETE("/users/profiles/:id", userProfileHandler.DeleteUserProfile)
	router.GET("/users/profiles/me", userProfileHandler.GetAuthenticatedUserProfile)

	return router
}

func TestCreateUserProfile(t *testing.T) {
	testCases := []struct {
		name              string
		userProfileCreate models.UserProfileCreate
		mockServiceReturn *models.UserProfileDTO
		mockServiceError  error
		expectedStatus    int
		expectedMessage   string
	}{
		{
			name: "Successful profile creation with all fields",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "John Doe",
				Bio:      "Test bio",
			},
			mockServiceReturn: &models.UserProfileDTO{
				ID:       uuid.New(),
				Fullname: "John Doe",
				Bio:      "Test bio",
			},
			mockServiceError: nil,
			expectedStatus:   http.StatusCreated,
			expectedMessage:  "User profile created successfully",
		},
		{
			name: "Successful profile creation with minimal fields",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "Jane Smith",
			},
			mockServiceReturn: &models.UserProfileDTO{
				ID:       uuid.New(),
				Fullname: "Jane Smith",
			},
			mockServiceError: nil,
			expectedStatus:   http.StatusCreated,
			expectedMessage:  "User profile created successfully",
		},
		{
			name: "Profile creation with empty fullname",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "",
				Bio:      "Test bio",
			},
			mockServiceReturn: nil,
			mockServiceError:  fmt.Errorf("fullname cannot be empty"),
			expectedStatus:    http.StatusBadRequest,
			expectedMessage:   "Validation failed",
		},
		{
			name: "Profile creation with very long fullname",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "John Jacob Jingleheimer Schmidt The Third Of His Name Long Lastname Extended Version",
				Bio:      "Very long name test bio",
			},
			mockServiceReturn: &models.UserProfileDTO{
				ID:       uuid.New(),
				Fullname: "John Jacob Jingleheimer Schmidt The Third Of His Name Long Lastname Extended Version",
				Bio:      "Very long name test bio",
			},
			mockServiceError: nil,
			expectedStatus:   http.StatusCreated,
			expectedMessage:  "User profile created successfully",
		},
		{
			name: "Profile creation with special characters in fullname",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "John-Doe Jr. & Associates",
				Bio:      "Special characters test bio",
			},
			mockServiceReturn: &models.UserProfileDTO{
				ID:       uuid.New(),
				Fullname: "John-Doe Jr. & Associates",
				Bio:      "Special characters test bio",
			},
			mockServiceError: nil,
			expectedStatus:   http.StatusCreated,
			expectedMessage:  "User profile created successfully",
		},
		{
			name: "Profile creation with service error",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "Jane Doe",
				Bio:      "Another test bio",
			},
			mockServiceReturn: nil,
			mockServiceError:  fmt.Errorf("internal server error"),
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "Failed to create user profile",
		},
		{
			name: "Profile creation with maximum length bio",
			userProfileCreate: models.UserProfileCreate{
				UserID:   uuid.New(),
				Fullname: "Max Bio User",
				Bio:      "This is a very long bio that tests the maximum length of a bio field which could potentially be quite extensive and contain multiple sentences to ensure that longer biographical information can be properly handled by the system without any truncation or unexpected behavior.",
			},
			mockServiceReturn: &models.UserProfileDTO{
				ID:       uuid.New(),
				Fullname: "Max Bio User",
				Bio:      "This is a very long bio that tests the maximum length of a bio field which could potentially be quite extensive and contain multiple sentences to ensure that longer biographical information can be properly handled by the system without any truncation or unexpected behavior.",
			},
			mockServiceError: nil,
			expectedStatus:   http.StatusCreated,
			expectedMessage:  "User profile created successfully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			// Prepare request body
			jsonBody, _ := json.Marshal(tc.userProfileCreate)
			req, _ := http.NewRequest("POST", "/users/profiles", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Mock service response
			if tc.expectedStatus != http.StatusBadRequest {
				mockUserProfileService.
					On("CreateProfile", mock.Anything, mock.MatchedBy(func(p *models.UserProfileCreate) bool {
						return p.Fullname == tc.userProfileCreate.Fullname &&
							p.Bio == tc.userProfileCreate.Bio
					})).
					Return(tc.mockServiceReturn, tc.mockServiceError)
			}

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Contains(t, response, "message")
			assert.Equal(t, tc.expectedMessage, response["message"])

			if tc.mockServiceReturn != nil {
				assert.Contains(t, response, "data")
			}

			// Add logging
			mockLogger.Info(fmt.Sprintf("Test case: %s, Status: %d", tc.name, w.Code))

			mockUserProfileService.AssertExpectations(t)
		})
	}
}

func TestListUserProfiles(t *testing.T) {
	testCases := []struct {
		name             string
		limit            int
		offset           int
		sortBy           string
		sortOrder        string
		expectedStatus   int
		expectedProfiles []models.UserProfileDTO
		serviceError     error
		expectedErrorMsg string
	}{
		{
			name:           "Successful list with default parameters",
			limit:          10,
			offset:         0,
			sortBy:         "created_at",
			sortOrder:      "desc",
			expectedStatus: http.StatusOK,
			expectedProfiles: []models.UserProfileDTO{
				{
					ID:       uuid.New(),
					Fullname: "John Doe",
					Bio:      "Test bio 1",
				},
				{
					ID:       uuid.New(),
					Fullname: "Jane Doe",
					Bio:      "Test bio 2",
				},
			},
		},
		{
			name:             "List with custom pagination",
			limit:            5,
			offset:           10,
			sortBy:           "fullname",
			sortOrder:        "asc",
			expectedStatus:   http.StatusOK,
			expectedProfiles: []models.UserProfileDTO{},
		},
		{
			name:             "Service returns error",
			limit:            10,
			offset:           0,
			sortBy:           "created_at",
			sortOrder:        "desc",
			expectedStatus:   http.StatusInternalServerError,
			serviceError:     errors.New("database error"),
			expectedErrorMsg: "Failed to retrieve user profiles",
		},
		{
			name:             "Invalid sort order fallback to default",
			limit:            10,
			offset:           0,
			sortBy:           "created_at",
			sortOrder:        "invalid",
			expectedStatus:   http.StatusOK,
			expectedProfiles: []models.UserProfileDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			if tc.sortOrder == "asc" {
				tc.sortOrder = base.SortAscending
			} else {
				tc.sortOrder = base.SortDescending
			}

			// Prepare list options
			listOptions := base.ListOptions{
				Page:      tc.offset/tc.limit + 1,
				PerPage:   tc.limit,
				SortBy:    tc.sortBy,
				SortOrder: tc.sortOrder,
			}

			// Mock service response
			if tc.serviceError != nil {
				mockUserProfileService.
					On("ListProfiles", mock.Anything, listOptions).
					Return([]models.UserProfileDTO(nil), 0, tc.serviceError)
			} else {
				mockUserProfileService.
					On("ListProfiles", mock.Anything, listOptions).
					Return(tc.expectedProfiles, len(tc.expectedProfiles), nil)
			}

			// Perform request
			url := fmt.Sprintf("/users/profiles?limit=%d&offset=%d&sort_by=%s&sort_order=%s",
				tc.limit, tc.offset, tc.sortBy, tc.sortOrder)
			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Enhanced logging for better visibility
			mockLogger.Info(fmt.Sprintf("Test Case: %s", tc.name))
			mockLogger.Info(fmt.Sprintf("Request URL: %s", url))
			mockLogger.Info(fmt.Sprintf("Expected Status: %d", tc.expectedStatus))
			mockLogger.Info(fmt.Sprintf("Actual Status: %d", w.Code))

			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "pagination")
				assert.Contains(t, response, "message")
				assert.Equal(t, "Data retrieved successfully", response["message"])

				mockLogger.Info("Response contains expected data and pagination")
			} else {
				assert.Contains(t, response, "message")
				if tc.expectedErrorMsg != "" {
					assert.Contains(t, response["message"], tc.expectedErrorMsg)
					mockLogger.Info(fmt.Sprintf("Error message contains: %s", tc.expectedErrorMsg))
				}
			}

			// Add logging
			mockLogger.Info(fmt.Sprintf("Test case: %s, Status: %d", tc.name, w.Code))

			mockUserProfileService.AssertExpectations(t)
		})
	}
}

func TestUpdateUserProfile(t *testing.T) {
	testCases := []struct {
		name               string
		profileID          uuid.UUID
		updateProfile      models.UserProfileUpdate
		mockServiceReturn  *models.UserProfileDTO
		mockServiceError   error
		expectedStatus     int
		expectedMessage    string
		invalidRequestBody bool
	}{
		{
			name:      "Successful profile update with all fields",
			profileID: uuid.New(),
			updateProfile: models.UserProfileUpdate{
				Fullname: "Updated Name",
				Bio:      "Updated comprehensive bio",
			},
			mockServiceReturn: &models.UserProfileDTO{
				Fullname: "Updated Name",
				Bio:      "Updated comprehensive bio",
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "User profile updated successfully",
		},
		{
			name:      "Successful profile update with minimal fields",
			profileID: uuid.New(),
			updateProfile: models.UserProfileUpdate{
				Fullname: "Minimal Update",
			},
			mockServiceReturn: &models.UserProfileDTO{
				Fullname: "Minimal Update",
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "User profile updated successfully",
		},
		{
			name:      "Profile update with service error",
			profileID: uuid.New(),
			updateProfile: models.UserProfileUpdate{
				Fullname: "Error Update",
			},
			mockServiceError: errors.New("database error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedMessage:  "Failed to update user profile",
		},
		{
			name:               "Invalid request body",
			profileID:          uuid.New(),
			invalidRequestBody: true,
			expectedStatus:     http.StatusBadRequest,
			expectedMessage:    "Invalid request body",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			// Prepare request body
			var jsonBody []byte
			var req *http.Request

			if tc.invalidRequestBody {
				// Send invalid JSON
				jsonBody = []byte("{invalid json}")
			} else {
				tc.updateProfile.ID = tc.profileID
				jsonBody, _ = json.Marshal(tc.updateProfile)
			}

			req, _ = http.NewRequest("PUT", "/users/profiles/"+tc.profileID.String(), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Setup mock expectations
			if !tc.invalidRequestBody {
				if tc.mockServiceError != nil {
					mockUserProfileService.On("UpdateProfile", mock.Anything, &tc.updateProfile).Return((*models.UserProfileDTO)(nil), tc.mockServiceError)
				} else {
					mockUserProfileService.On("UpdateProfile", mock.Anything, mock.MatchedBy(func(up *models.UserProfileUpdate) bool {
						return up.Fullname == tc.updateProfile.Fullname && up.Bio == tc.updateProfile.Bio
					})).
						Return(tc.mockServiceReturn, nil)
				}
			}

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Contains(t, response, "message")
			assert.Equal(t, tc.expectedMessage, response["message"])

			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "data")
			}

			mockUserProfileService.AssertExpectations(t)
		})
	}
}

func TestUpdateUserAvatar(t *testing.T) {
	testCases := []struct {
		name               string
		profileID          uuid.UUID
		avatarFileName     string
		avatarContent      []byte
		mockServiceReturn  *models.UserProfileDTO
		mockServiceError   error
		expectedStatus     int
		expectedMessage    string
		invalidRequestBody bool
	}{
		{
			name:           "Successful avatar update",
			profileID:      uuid.New(),
			avatarFileName: "avatar.jpg",
			avatarContent:  []byte("fake image content"),
			mockServiceReturn: &models.UserProfileDTO{
				ID:        uuid.New(),
				AvatarUrl: "https://example.com/avatar.jpg",
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "User profile avatar updated successfully",
		},
		{
			name:             "Service returns error",
			profileID:        uuid.New(),
			avatarFileName:   "avatar.jpg",
			avatarContent:    []byte("fake image content"),
			mockServiceError: errors.New("avatar update failed"),
			expectedStatus:   http.StatusInternalServerError,
			expectedMessage:  "Failed to update user profile avatar",
		},
		{
			name:               "Invalid request body",
			profileID:          uuid.New(),
			invalidRequestBody: true,
			expectedStatus:     http.StatusBadRequest,
			expectedMessage:    "Invalid multipart form request",
		},
		{
			name:            "No avatar file uploaded",
			profileID:       uuid.New(),
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Avatar file is required",
		},
		{
			name:            "Large file upload",
			profileID:       uuid.New(),
			avatarFileName:  "large_avatar.jpg",
			avatarContent:   make([]byte, 10*1024*1024), // 10MB file
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "File size exceeds maximum limit of 2MB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			// Prepare request
			var req *http.Request
			var body *bytes.Buffer
			var writer *multipart.Writer

			if !tc.invalidRequestBody {
				body = &bytes.Buffer{}
				writer = multipart.NewWriter(body)

				if len(tc.avatarContent) > 0 {
					part, _ := writer.CreateFormFile("avatar", tc.avatarFileName)
					part.Write(tc.avatarContent)
				}
				writer.Close()

				req, _ = http.NewRequest("PUT", "/users/profiles/"+tc.profileID.String()+"/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				// Invalid request body scenario
				body = bytes.NewBufferString("invalid body")
				req, _ = http.NewRequest("PUT", "/users/profiles/"+tc.profileID.String()+"/avatar", body)
				req.Header.Set("Content-Type", "application/json")
			}

			if tc.name != "Large file upload" {
				// Setup mock expectations
				if tc.mockServiceError != nil {
					mockUserProfileService.On("UpdateProfileAvatar", mock.Anything, mock.Anything).
						Return((*models.UserProfileDTO)(nil), tc.mockServiceError)
				} else if !tc.invalidRequestBody && len(tc.avatarContent) > 0 {
					mockUserProfileService.On("UpdateProfileAvatar", mock.Anything, mock.Anything).
						Return(tc.mockServiceReturn, nil)
				}
			}

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Contains(t, response, "message")
			assert.Equal(t, tc.expectedMessage, response["message"])

			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "data")
			}

			// Log test details
			mockLogger.Info(fmt.Sprintf("Test Case: %s", tc.name))
			mockLogger.Info(fmt.Sprintf("Expected Status: %d", tc.expectedStatus))
			mockLogger.Info(fmt.Sprintf("Actual Status: %d", w.Code))

			mockUserProfileService.AssertExpectations(t)
		})
	}
}

func TestDeleteUserProfile(t *testing.T) {
	testCases := []struct {
		name             string
		profileID        uuid.UUID
		mockServiceError error
		expectedStatus   int
		expectedMessage  string
	}{
		{
			name:             "Successful Delete",
			profileID:        uuid.New(),
			mockServiceError: nil,
			expectedStatus:   http.StatusOK,
			expectedMessage:  "User profile deleted successfully",
		},
		{
			name:             "Profile Not Found",
			profileID:        uuid.New(),
			mockServiceError: errLib.New(errLib.ErrNotFound, "Profile not found", nil),
			expectedStatus:   http.StatusNotFound,
			expectedMessage:  "Profile not found",
		},
		{
			name:             "Unauthorized Delete",
			profileID:        uuid.New(),
			mockServiceError: errLib.New(errLib.ErrUnauthorized, "Unauthorized to delete profile", nil),
			expectedStatus:   http.StatusUnauthorized,
			expectedMessage:  "Unauthorized to delete profile",
		},
		{
			name:             "Internal Server Error",
			profileID:        uuid.New(),
			mockServiceError: errLib.New(errLib.ErrInternal, "Failed to delete user profile", nil),
			expectedStatus:   http.StatusInternalServerError,
			expectedMessage:  "Failed to delete user profile",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			// Setup mock expectations
			if tc.mockServiceError != nil {
				mockUserProfileService.On("DeleteProfile", mock.Anything, tc.profileID.String()).
					Return(tc.mockServiceError)
			} else {
				mockUserProfileService.On("DeleteProfile", mock.Anything, tc.profileID.String()).
					Return(nil)
			}

			// Perform request
			req, _ := http.NewRequest("DELETE", "/users/profiles/"+tc.profileID.String(), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Contains(t, response, "message")
			assert.Equal(t, tc.expectedMessage, response["message"])

			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "data")
			}

			// Log test details
			mockLogger.Info(fmt.Sprintf("Test Case: %s", tc.name))
			mockLogger.Info(fmt.Sprintf("Expected Status: %d", tc.expectedStatus))
			mockLogger.Info(fmt.Sprintf("Actual Status: %d", w.Code))

			mockUserProfileService.AssertExpectations(t)
		})
	}
}

func TestGetUserProfileByID(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserProfileService := new(MockUserProfileService)
	userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

	router := setupUserProfileTestRouter(userProfileHandler)

	profileID := uuid.New()

	// Mock service response
	mockUserProfileService.On("GetProfileByID", mock.Anything, profileID.String()).Return(&models.UserProfileDTO{
		ID:       profileID,
		Fullname: "John Doe",
		Bio:      "Test bio",
	}, nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/profiles/"+profileID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User profile retrieved successfully", response["message"])

	mockUserProfileService.AssertExpectations(t)
}

func TestGetAuthenticatedUserProfile(t *testing.T) {
	testCases := []struct {
		name             string
		mockServiceError error
		expectedStatus   int
		expectedMessage  string
		userID           string
		userEmail        string
		userRole         string
		userBadges       []string
	}{
		{
			name:             "Successful retrieval of authenticated user profile",
			mockServiceError: nil,
			expectedStatus:   http.StatusOK,
			expectedMessage:  "User profile retrieved successfully",
			userID:           "test-user-id",
			userEmail:        "test@example.com",
			userRole:         "user",
			userBadges:       []string{"explorer"},
		},
		{
			name:             "Error retrieving user profile - Not Found",
			mockServiceError: errLib.New(errLib.ErrNotFound, "Profile not found", nil),
			expectedStatus:   http.StatusNotFound,
			expectedMessage:  "Profile not found",
			userID:           "test-user-id",
			userEmail:        "test@example.com",
			userRole:         "user",
			userBadges:       []string{},
		},
		{
			name:             "Error retrieving user profile - Unauthorized",
			mockServiceError: errLib.New(errLib.ErrUnauthorized, "Unauthorized to access profile", nil),
			expectedStatus:   http.StatusUnauthorized,
			expectedMessage:  "Unauthorized to access profile",
			userID:           "test-user-id",
			userEmail:        "test@example.com",
			userRole:         "user",
			userBadges:       []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := logger.NewLogger(logger.LoggerConfig{})
			mockUserProfileService := new(MockUserProfileService)
			userProfileHandler := handlers.NewUserProfileHandler(mockUserProfileService, mockLogger)

			router := setupUserProfileTestRouter(userProfileHandler)

			// Setup mock expectations
			if tc.mockServiceError == nil {
				mockUserProfileService.On("GetProfileByUserID", mock.Anything, mock.AnythingOfType("string")).
					Return(&models.UserProfileDTO{
						ID:       uuid.New(),
						Fullname: "John Doe",
						Bio:      "Test bio",
					}, nil)
			} else {
				mockUserProfileService.On("GetProfileByUserID", mock.Anything, mock.AnythingOfType("string")).Return((*models.UserProfileDTO)(nil), tc.mockServiceError)
			}

			// Perform request
			req, _ := http.NewRequest("GET", "/users/profiles/me", nil)

			// Simulate middleware context with user details
			req = req.WithContext(
				context.WithValue(req.Context(), "user_id", tc.userID),
			)
			req = req.WithContext(
				context.WithValue(req.Context(), "email", tc.userEmail),
			)
			req = req.WithContext(
				context.WithValue(req.Context(), "role", tc.userRole),
			)
			req = req.WithContext(
				context.WithValue(req.Context(), "badges", tc.userBadges),
			)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check response based on expected status
			switch tc.expectedStatus {
			case http.StatusOK:
				assert.Contains(t, response, "data")
				assert.Equal(t, tc.expectedMessage, response["message"])
			default:
				assert.Contains(t, response, "error")
				assert.Equal(t, tc.expectedMessage, response["message"])
			}

			mockUserProfileService.AssertExpectations(t)

			// Log test details
			mockLogger.Info(fmt.Sprintf("Test Case: %s", tc.name))
			mockLogger.Info(fmt.Sprintf("Expected Status: %d", tc.expectedStatus))
			mockLogger.Info(fmt.Sprintf("Actual Status: %d", w.Code))
		})
	}
}
