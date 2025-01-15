package server

import (
	"errors"
	"fmt"
	"net/http"
	"sports-news-api/internal/app/domain"
	"sports-news-api/internal/app/transport"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s Server) GetArticleById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["article_id"]
	clientId := r.URL.Query().Get("clientId")

	article, err := s.articlesService.GetArticleById(r.Context(), articleId, clientId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			RespondWithFailure("Article with such id is missing", http.StatusBadRequest, w, r)
		} else {
			RespondWithError(fmt.Sprintf("error getting article by id: %v", err), http.StatusInternalServerError, w, r)
		}
		return
	}

	responseArticle := transport.Article{
		Id:              article.ID(),
		ClientId:        article.ClientID(),
		ClientArticleId: article.ClientArticleId(),
		Title:           article.Title(),
		Content:         article.Content(),
		PublishDate:     article.PublishDate().Format(time.RFC3339),
	}

	response := transport.ArticleResponse{
		Status: StatusSuccess,
		Data:   responseArticle,
		Metadata: transport.Metadata{
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	RespondWithSuccess(response, w, r)
}

func (s Server) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	offset := r.URL.Query().Get("offset")
	var defaultQueryLimit int64 = 100

	var offsetInt int64 = 0
	if offset != "" {
		var err error
		offsetInt, err = strconv.ParseInt(offset, 10, 64)
		if err != nil {
			RespondWithFailure(fmt.Sprintf("invalid offset value: %v", err), http.StatusBadRequest, w, r)
			return
		}
	}
	articles, err := s.articlesService.GetAllArticles(r.Context(), clientId, defaultQueryLimit, offsetInt)
	if err != nil {
		RespondWithError(fmt.Sprintf("error getting articles: %v", err), http.StatusInternalServerError, w, r)
		return
	}

	var responseArticles []transport.Article

	for _, article := range articles {
		responseArticle := transport.Article{
			Id:              article.ID(),
			ClientId:        article.ClientID(),
			ClientArticleId: article.ClientArticleId(),
			Title:           article.Title(),
			Content:         article.Content(),
			PublishDate:     article.PublishDate().Format(time.RFC3339),
		}
		responseArticles = append(responseArticles, responseArticle)
	}

	response := transport.ArticlesResponse{
		Status: StatusSuccess,
		Data:   responseArticles,
		Metadata: transport.Metadata{
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	RespondWithSuccess(response, w, r)
}
