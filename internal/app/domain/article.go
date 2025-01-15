package domain

import (
	"fmt"
	"time"
)

// Article is a domain news article.
type Article struct {
	id              string
	clientId        string
	clientArticleId string
	title           string
	content         string
	lastModified    int64
	publishDate     time.Time
}

func NewArticle(id, clientId, clientArticleId, title, content string, lastModified int64, publishDate time.Time) (*Article, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: article id is required", ErrRequired)
	}

	if clientId == "" {
		return nil, fmt.Errorf("%w: clientId is required", ErrRequired)
	}

	if clientArticleId == "" {
		return nil, fmt.Errorf("%w: clientArticleId is required", ErrRequired)
	}

	if title == "" {
		return nil, fmt.Errorf("%w: article title is required", ErrRequired)
	}

	if content == "" {
		return nil, fmt.Errorf("%w: article content is required", ErrRequired)
	}

	return &Article{
		id:              id,
		clientId:        clientId,
		clientArticleId: clientArticleId,
		title:           title,
		content:         content,
		lastModified:    lastModified,
		publishDate:     publishDate,
	}, nil
}

// ID returns the article id.
func (a *Article) ID() string {
	return a.id
}

// Client Id returns the article id.
func (a *Article) ClientID() string {
	return a.clientId
}

// Client Article Id returns the article id.
func (a *Article) ClientArticleId() string {
	return a.clientArticleId
}

// Title returns the article title.
func (a *Article) Title() string {
	return a.title
}

// Content returns the article content.
func (a *Article) Content() string {
	return a.content
}

// PublishDate returns the article publish date.
func (a *Article) LastModified() int64 {
	return a.lastModified
}

// PublishDate returns the article publish date.
func (a *Article) PublishDate() time.Time {
	return a.publishDate
}
