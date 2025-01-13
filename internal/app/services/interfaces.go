package services

import (
	"context"
	"sports-news-api/internal/app/models"
)

type ArticlesRepo interface {
	UpsertEcbArticle(ctx context.Context, getEcbArticleResponse models.EcbArticle) error
	GetEcbArticleById(ctx context.Context, id string) (interface{}, error)
	GetAllEcbArticles(ctx context.Context, limit int64) ([]interface{}, error)
}
