package server

import (
	"fmt"
	"net/http"
	"sports-news-api/internal/app/utils"

	"github.com/gorilla/mux"
)

func (s Server) GetAllArticles(w http.ResponseWriter, r *http.Request) {

	result, err := s.articlesService.GetAllEcbArticles(r.Context(), 100)
	if err != nil {
		RespondWithError(fmt.Sprintf("Error getting articles: %v", err), http.StatusInternalServerError, w, r)
		return
	}

	RespondWithSuccess(result, w, r)
}

func (s Server) GetArticleById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId := vars["article_id"]

	result, err := s.articlesService.GetEcbArticleById(r.Context(), articleId)
	if err != nil {
		if err == utils.ErrNotFound {
			RespondWithFailure("Article with such id is missing", http.StatusBadRequest, w, r)
		} else {
			RespondWithError(fmt.Sprintf("Error getting article by id: %v", err), http.StatusInternalServerError, w, r)
		}

		return
	}

	RespondWithSuccess(result, w, r)
}
