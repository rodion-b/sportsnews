package models

import "time"

type GetEcbArticleResponse struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"accountId"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description any       `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Coordinates []float64 `json:"coordinates"`
	CommentsOn  bool      `json:"commentsOn"`
	Copyright   any       `json:"copyright"`
	PublishFrom int64     `json:"publishFrom"`
	PublishTo   int       `json:"publishTo"`
	Tags        []struct {
		ID    int    `json:"id"`
		Label string `json:"label"`
	} `json:"tags"`
	Platform       string `json:"platform"`
	Language       string `json:"language"`
	AdditionalInfo struct {
	} `json:"additionalInfo"`
	CanonicalURL string `json:"canonicalUrl"`
	References   []struct {
		Label any    `json:"label"`
		Sid   string `json:"sid"`
		ID    int    `json:"id"`
		Type  string `json:"type"`
	} `json:"references"`
	Related  []any `json:"related"`
	Metadata struct {
	} `json:"metadata"`
	TitleTranslations any    `json:"titleTranslations"`
	LastModified      int64  `json:"lastModified"`
	TitleURLSegment   string `json:"titleUrlSegment"`
	Body              any    `json:"body"`
	Author            any    `json:"author"`
	Subtitle          any    `json:"subtitle"`
	Summary           any    `json:"summary"`
	HotlinkURL        string `json:"hotlinkUrl"`
	Duration          int    `json:"duration"`
	ContentSummary    any    `json:"contentSummary"`
	LeadMedia         struct {
		ID             int       `json:"id"`
		AccountID      int       `json:"accountId"`
		Type           string    `json:"type"`
		Title          string    `json:"title"`
		Description    any       `json:"description"`
		Date           time.Time `json:"date"`
		Location       any       `json:"location"`
		Coordinates    []float64 `json:"coordinates"`
		CommentsOn     bool      `json:"commentsOn"`
		Copyright      any       `json:"copyright"`
		PublishFrom    int64     `json:"publishFrom"`
		PublishTo      int       `json:"publishTo"`
		Tags           []any     `json:"tags"`
		Platform       string    `json:"platform"`
		Language       string    `json:"language"`
		AdditionalInfo struct {
		} `json:"additionalInfo"`
		CanonicalURL string `json:"canonicalUrl"`
		References   []any  `json:"references"`
		Related      []any  `json:"related"`
		Metadata     struct {
		} `json:"metadata"`
		TitleTranslations any    `json:"titleTranslations"`
		LastModified      int64  `json:"lastModified"`
		TitleURLSegment   string `json:"titleUrlSegment"`
		Subtitle          any    `json:"subtitle"`
		Variants          []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
			Tag    struct {
				ID    int    `json:"id"`
				Label string `json:"label"`
			} `json:"tag"`
		} `json:"variants"`
		OnDemandURL     string `json:"onDemandUrl"`
		OriginalDetails struct {
			Width       int     `json:"width"`
			Height      int     `json:"height"`
			AspectRatio float64 `json:"aspectRatio"`
		} `json:"originalDetails"`
		ImageURL string `json:"imageUrl"`
	} `json:"leadMedia"`
	ImageURL    string `json:"imageUrl"`
	OnDemandURL string `json:"onDemandUrl"`
}
