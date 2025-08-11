package services_test

import (
	"context"
	"testing"

	"github.com/holycann/cultour-backend/internal/gemini/models"
	geminiServices "github.com/holycann/cultour-backend/internal/gemini/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeminiClient is a mock implementation of GeminiClient
type MockGeminiClient struct {
	mock.Mock
}

func (m *MockGeminiClient) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func (m *MockGeminiClient) GenerateRuleBasedResponse(ctx context.Context, prompt string, rules []models.Rule) (string, error) {
	args := m.Called(ctx, prompt, rules)
	return args.String(0), args.Error(1)
}

func (m *MockGeminiClient) AnalyzeText(ctx context.Context, text string) (*models.TextAnalysis, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(*models.TextAnalysis), args.Error(1)
}

func (m *MockGeminiClient) TranslateText(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error) {
	args := m.Called(ctx, text, sourceLanguage, targetLanguage)
	return args.String(0), args.Error(1)
}

func (m *MockGeminiClient) SummarizeText(ctx context.Context, text string, maxLength int) (string, error) {
	args := m.Called(ctx, text, maxLength)
	return args.String(0), args.Error(1)
}

func (m *MockGeminiClient) DetectLanguage(ctx context.Context, text string) (string, error) {
	args := m.Called(ctx, text)
	return args.String(0), args.Error(1)
}

func TestGenerateResponse(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	prompt := "Tell me about the history of Jakarta"
	expectedResponse := "Jakarta, formerly known as Batavia, is the capital city of Indonesia..."

	mockClient.On("GenerateResponse", mock.Anything, prompt).Return(expectedResponse, nil)

	result, err := geminiService.GenerateResponse(context.Background(), prompt)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)

	mockClient.AssertExpectations(t)
}

func TestGenerateRuleBasedResponse(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	prompt := "Describe a cultural event in Indonesia"
	rules := []models.Rule{
		{
			Name:        "Length",
			Description: "Response should be between 50-100 words",
		},
		{
			Name:        "Tone",
			Description: "Use an informative and engaging tone",
		},
	}
	expectedResponse := "The Bali Arts Festival is a vibrant celebration of Balinese culture..."

	mockClient.On("GenerateRuleBasedResponse", mock.Anything, prompt, rules).Return(expectedResponse, nil)

	result, err := geminiService.GenerateRuleBasedResponse(context.Background(), prompt, rules)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)

	mockClient.AssertExpectations(t)
}

func TestAnalyzeText(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	text := "Jakarta is the capital of Indonesia, known for its rich cultural heritage and modern development."
	expectedAnalysis := &models.TextAnalysis{
		Sentiment:     "Neutral",
		KeyTopics:     []string{"Jakarta", "Capital", "Culture", "Development"},
		Complexity:    "Medium",
		LanguageStyle: "Informative",
	}

	mockClient.On("AnalyzeText", mock.Anything, text).Return(expectedAnalysis, nil)

	result, err := geminiService.AnalyzeText(context.Background(), text)

	assert.NoError(t, err)
	assert.Equal(t, expectedAnalysis, result)

	mockClient.AssertExpectations(t)
}

func TestTranslateText(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	text := "Selamat datang di Jakarta"
	sourceLanguage := "id"
	targetLanguage := "en"
	expectedTranslation := "Welcome to Jakarta"

	mockClient.On("TranslateText", mock.Anything, text, sourceLanguage, targetLanguage).Return(expectedTranslation, nil)

	result, err := geminiService.TranslateText(context.Background(), text, sourceLanguage, targetLanguage)

	assert.NoError(t, err)
	assert.Equal(t, expectedTranslation, result)

	mockClient.AssertExpectations(t)
}

func TestSummarizeText(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	text := "Jakarta is the capital and largest city of Indonesia. It is located on the northwest coast of Java, the world's most populous island. As a major economic, cultural, and political center, Jakarta plays a crucial role in Indonesia's development. The city is known for its diverse population, rich history, and rapid urbanization."
	maxLength := 100
	expectedSummary := "Jakarta, Indonesia's capital on Java's northwest coast, is a major economic, cultural, and political center known for its diversity and rapid urbanization."

	mockClient.On("SummarizeText", mock.Anything, text, maxLength).Return(expectedSummary, nil)

	result, err := geminiService.SummarizeText(context.Background(), text, maxLength)

	assert.NoError(t, err)
	assert.Equal(t, expectedSummary, result)

	mockClient.AssertExpectations(t)
}

func TestDetectLanguage(t *testing.T) {
	mockClient := new(MockGeminiClient)
	geminiService := geminiServices.NewGeminiService(mockClient)

	text := "Selamat datang di Jakarta"
	expectedLanguage := "Indonesian"

	mockClient.On("DetectLanguage", mock.Anything, text).Return(expectedLanguage, nil)

	result, err := geminiService.DetectLanguage(context.Background(), text)

	assert.NoError(t, err)
	assert.Equal(t, expectedLanguage, result)

	mockClient.AssertExpectations(t)
}
