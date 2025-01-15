package models

import "time"

type GetEcbArticlesResponse struct {
	PageInfo struct {
		Page       int `json:"page"`
		NumPages   int `json:"numPages"`
		PageSize   int `json:"pageSize"`
		NumEntries int `json:"numEntries"`
	} `json:"pageInfo"`
	Content []struct {
		ID           int       `json:"id"`
		AccountID    int       `json:"accountId"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		Date         time.Time `json:"date"`
		LastModified int64     `json:"lastModified"`
	} `json:"content"`
}
