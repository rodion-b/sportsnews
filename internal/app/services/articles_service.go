package services

import (
	"context"
	"sports-news-api/internal/app/models"
)

type ArticlesService struct {
	articlesRepo ArticlesRepo
}

func NewArticlesService(articlesRepo ArticlesRepo) ArticlesService {
	return ArticlesService{articlesRepo: articlesRepo}
}

func (s ArticlesService) UpsertEcbArticle(ctx context.Context, ecbArticle models.EcbArticle) error {
	return s.articlesRepo.UpsertEcbArticle(ctx, ecbArticle)
}

func (s ArticlesService) GetEcbArticleById(ctx context.Context, id string) (interface{}, error) {
	return s.articlesRepo.GetEcbArticleById(ctx, id)
}

func (s ArticlesService) GetAllEcbArticles(ctx context.Context, limit int64) ([]interface{}, error) {
	return s.articlesRepo.GetAllEcbArticles(ctx, limit)
}
