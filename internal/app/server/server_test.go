package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of ArticlesService interface
type MockArticlesService struct {
	mock.Mock
}

func (m *MockArticlesService) GetEcbArticleById(ctx context.Context, id string) (interface{}, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

func (m *MockArticlesService) GetAllEcbArticles(ctx context.Context, limit int64) ([]interface{}, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]interface{}), args.Error(1)
}

func newTestServer() Server {
	return Server{
		articlesService: &MockArticlesService{},
	}
}

func TestGetAllArticles_Success(t *testing.T) {
	mockService := &MockArticlesService{}
	server := Server{articlesService: mockService}

	articles := []interface{}{"Article 1", "Article 2"}
	mockService.On("GetAllEcbArticles", mock.Anything, int64(100)).Return(articles, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()

	server.GetAllArticles(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response SuccessResponse
	_ = json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, articles, response.Data)
}

func TestGetAllArticles_Error(t *testing.T) {
	mockService := &MockArticlesService{}
	server := Server{articlesService: mockService}

	mockService.On("GetAllEcbArticles", mock.Anything, int64(100)).Return([]interface{}{}, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()

	server.GetAllArticles(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response ErrorResponse
	_ = json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "Error getting articles")
}

func TestGetArticleById_Success(t *testing.T) {
	mockService := &MockArticlesService{}
	server := Server{articlesService: mockService}

	articleID := "123"
	expectedArticle := map[string]interface{}{
		"id":      "123",
		"title":   "Sample Article",
		"content": "This is a sample article",
	}
	mockService.On("GetEcbArticleById", mock.Anything, articleID).Return(expectedArticle, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles/123", nil)
	req = mux.SetURLVars(req, map[string]string{"article_id": articleID})
	rec := httptest.NewRecorder()

	server.GetArticleById(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response SuccessResponse
	_ = json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, expectedArticle, response.Data)
}

func TestGetArticleById_Error(t *testing.T) {
	mockService := &MockArticlesService{}
	server := Server{articlesService: mockService}

	articleID := "123"
	mockService.On("GetEcbArticleById", mock.Anything, articleID).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/articles/123", nil)
	req = mux.SetURLVars(req, map[string]string{"article_id": articleID})
	rec := httptest.NewRecorder()

	server.GetArticleById(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response ErrorResponse
	_ = json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "Error getting article by id")
}
