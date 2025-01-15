package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sports-news-api/internal/app/domain"
	"sports-news-api/internal/app/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EcbFeedsService struct {
	httpClient      *http.Client
	baseUrl         *url.URL
	articlesService ArticlesService
}

func NewEcbFeedsService(host string, path string, httpTimeout time.Duration, articlesService ArticlesService) EcbFeedsService {
	client := &http.Client{Timeout: httpTimeout}

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}
	return EcbFeedsService{
		client,
		baseUrl,
		articlesService,
	}
}

func (s *EcbFeedsService) PollEcbArticles(ctx context.Context, pollSize int, refreshFrequency time.Duration) {
	ticker := time.NewTicker(refreshFrequency)
	for {
		select {
		case <-ticker.C:
			//get latest N articles
			ecbArticles, err := s.GetEcbArticlesIds(pollSize)
			if err != nil {
				log.Err(err).Msg("Unable to GetEcbArticles")
				continue //skipping the rest
			}
			for _, articleId := range ecbArticles {
				go func() {
					article, err := s.GetArticleById(articleId)
					if err != nil {
						log.Err(err).Msg("Unable to GetArticleById")
						return
					}
					//upserting article into db
					err = s.articlesService.UpsertArticle(ctx, *article)
					if err != nil {
						log.Err(err).Msg("Unable to UpsertEcbArticle")
						return
					}
				}()
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *EcbFeedsService) GetArticleById(id string) (*domain.Article, error) {
	//constructing article url
	baseUrl, err := url.Parse(s.baseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	// Append new path segment
	baseUrl.Path = path.Join(baseUrl.Path, id)
	finalUrl := baseUrl.String() // Full URL with query parameters

	resp, err := s.httpClient.Get(finalUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var ecbArticle models.GetEcbArticleResponse
	err = json.Unmarshal(body, &ecbArticle)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	domainArticle, err := domain.NewArticle(
		uuid.New().String(),
		domain.Ecb,
		strconv.Itoa(ecbArticle.ID),
		ecbArticle.Title,
		fmt.Sprintf("%v", ecbArticle.Body), //raw content
		ecbArticle.LastModified,
		ecbArticle.Date,
	)
	if err != nil {
		return nil, fmt.Errorf("error converting to domain article: %v", err)
	}

	return domainArticle, nil
}

/*
Gets Latest N articles from ECB
*/
func (s *EcbFeedsService) GetEcbArticlesIds(pageSize int) ([]string, error) {

	// Create a copy of the base URL to modify
	baseUrl, err := url.Parse(s.baseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	// Add query parameters
	query := baseUrl.Query()
	query.Set("pageSize", strconv.Itoa(pageSize))
	baseUrl.RawQuery = query.Encode()

	finalUrl := baseUrl.String() // Full URL with query parameters

	//Sending request
	resp, err := s.httpClient.Get(finalUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var getEcbArticlesResponse models.GetEcbArticlesResponse
	err = json.Unmarshal(body, &getEcbArticlesResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	//converting to ids only
	var articleIds []string
	for _, article := range getEcbArticlesResponse.Content {
		articleIds = append(articleIds, strconv.Itoa(article.ID))
	}
	return articleIds, nil
}
