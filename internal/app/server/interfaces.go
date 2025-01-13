package server

import (
	"context"
)

type ArticlesService interface {
	GetEcbArticleById(ctx context.Context, id string) (interface{}, error)
	GetAllEcbArticles(ctx context.Context, limit int64) ([]interface{}, error)
}
