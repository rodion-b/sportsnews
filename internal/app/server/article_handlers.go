package server

import (
	"fmt"
	"net/http"

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

	if articleId == "" {
		RespondWithFailure("Missing article_id in request", http.StatusBadRequest, w, r)
		return
	}

	result, err := s.articlesService.GetEcbArticleById(r.Context(), articleId)
	if err != nil {
		RespondWithError(fmt.Sprintf("Error getting article by id: %v", err), http.StatusInternalServerError, w, r)
		return
	}

	RespondWithSuccess(result, w, r)
}
