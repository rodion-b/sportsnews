package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sports-news-api/internal/app/models"
	"strconv"
	"time"
)

type EcbFeedsService struct {
	httpClient *http.Client
	baseUrl    *url.URL
}

func NewEcbFeedsService(host string, path string, httpTimeout time.Duration) EcbFeedsService {
	client := &http.Client{Timeout: httpTimeout}

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}
	return EcbFeedsService{
		client,
		baseUrl,
	}
}

func (s *EcbFeedsService) GetArticleById(id string) (*models.EcbArticle, error) {
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

	var ecbArticle models.EcbArticle
	err = json.Unmarshal(body, &ecbArticle)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &ecbArticle, nil
}

/*
Gets Latest N articles
*/
func (s *EcbFeedsService) GetEcbArticles(pageSize int) (*models.GetEcbArticlesResponse, error) {

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
	return &getEcbArticlesResponse, nil
}
