package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sports-news-api/internal/app/domain"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Mock ArticlesService implementing the domain.Article struct
type mockArticlesService struct{}

func (m *mockArticlesService) UpsertArticle(ctx context.Context, article domain.Article) error {
	return nil
}

func (m *mockArticlesService) GetArticleById(ctx context.Context, id string, clientId string) (*domain.Article, error) {
	if id == "valid-id" {
		article, _ := domain.NewArticle(
			"valid-id",
			clientId,
			"123",
			"Test Article",
			"This is a test article.",
			time.Now().Unix(),
			time.Now(),
		)
		return article, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockArticlesService) GetAllArticles(ctx context.Context, clientId string, limit int64, offset int64) ([]*domain.Article, error) {
	article1, _ := domain.NewArticle(
		"article-1",
		clientId,
		"101",
		"Article 1",
		"Content of article 1",
		time.Now().Unix(),
		time.Now(),
	)
	article2, _ := domain.NewArticle(
		"article-2",
		clientId,
		"102",
		"Article 2",
		"Content of article 2",
		time.Now().Unix(),
		time.Now(),
	)
	return []*domain.Article{article1, article2}, nil
}

// Test GetArticleById
func TestGetArticleById(t *testing.T) {
	service := &mockArticlesService{}
	server := NewServer(service)

	req := httptest.NewRequest("GET", "/article/valid-id?clientId=test-client", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/article/{article_id}", server.GetArticleById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), `"status":"success"`))
	assert.True(t, strings.Contains(rr.Body.String(), `"title":"Test Article"`))
}

// Test GetArticleById - Article Not Found
func TestGetArticleById_NotFound(t *testing.T) {
	service := &mockArticlesService{}
	server := NewServer(service)

	req := httptest.NewRequest("GET", "/article/invalid-id?clientId=test-client", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/article/{article_id}", server.GetArticleById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), `"status":"fail"`))
}

// Test GetAllArticles
func TestGetAllArticles(t *testing.T) {
	service := &mockArticlesService{}
	server := NewServer(service)

	req := httptest.NewRequest("GET", "/articles?clientId=test-client&offset=0", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/articles", server.GetAllArticles)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), `"status":"success"`))
	assert.True(t, strings.Contains(rr.Body.String(), `"title":"Article 1"`))
	assert.True(t, strings.Contains(rr.Body.String(), `"title":"Article 2"`))
}

// Test GetAllArticles - Invalid Offset
func TestGetAllArticles_InvalidOffset(t *testing.T) {
	service := &mockArticlesService{}
	server := NewServer(service)

	req := httptest.NewRequest("GET", "/articles?clientId=test-client&offset=invalid", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/articles", server.GetAllArticles)
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), `"status":"fail"`))
	assert.True(t, strings.Contains(rr.Body.String(), "invalid offset value"))
}
