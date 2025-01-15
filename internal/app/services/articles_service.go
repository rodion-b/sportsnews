package services

import (
	"context"
	"sports-news-api/internal/app/domain"
)

type ArticlesService struct {
	articlesRepo ArticlesRepo
}

func NewArticlesService(articlesRepo ArticlesRepo) ArticlesService {
	return ArticlesService{articlesRepo: articlesRepo}
}

func (s ArticlesService) UpsertArticle(ctx context.Context, article domain.Article) error {
	return s.articlesRepo.UpsertArticle(ctx, article)
}

func (s ArticlesService) GetArticleById(ctx context.Context, id string, clientId string) (*domain.Article, error) {
	return s.articlesRepo.GetArticleById(ctx, id, clientId)
}

func (s ArticlesService) GetAllArticles(ctx context.Context, clientId string, limit int64, offset int64) ([]*domain.Article, error) {
	return s.articlesRepo.GetAllArticles(ctx, clientId, limit, offset)
}
