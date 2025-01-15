package services

import (
	"context"
	"sports-news-api/internal/app/domain"
)

type ArticlesRepo interface {
	UpsertArticle(ctx context.Context, article domain.Article) error
	GetArticleById(ctx context.Context, id string, clientId string) (*domain.Article, error)
	GetAllArticles(ctx context.Context, clientId string, limit int64, offset int64) ([]*domain.Article, error)
}
