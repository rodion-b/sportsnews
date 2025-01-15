package repository

import (
	"time"
)

type Article struct {
	Id              string    `bson:"_id"`
	ClientId        string    `bson:"clientId"`
	ClientArticleId string    `bson:"clientArticleId"`
	Title           string    `bson:"title"`
	Description     string    `bson:"description"`
	Content         string    `bson:"content"`
	LastModified    int64     `bson:"lastModified"`
	PublishDate     time.Time `bson:"publishDate"`
}
