package repository

import (
	"context"
	"errors"
	"fmt"
	"sports-news-api/internal/app/domain"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticlesRepo struct {
	client *mongo.Client
	dbName string
}

func NewArticlesRepo(client *mongo.Client, dbName string) *ArticlesRepo {
	return &ArticlesRepo{client: client, dbName: dbName}
}

/*
If article already exists in the database we check if lastmodifed date is old and update the article
If article doesnt exist we insert new one
*/
func (r ArticlesRepo) UpsertArticle(ctx context.Context, article domain.Article) error {

	articlesCollection := r.client.Database(r.dbName).Collection(ArticlesCollection)

	var existingAritcle Article
	err := articlesCollection.FindOne(ctx, bson.M{"clientArticleId": article.ClientArticleId()}).Decode(&existingAritcle)

	//if not documents found insert one
	if errors.Is(err, mongo.ErrNoDocuments) {
		//inserting one
		repoArticle := bson.M{
			"_id":             article.ID(),
			"clientId":        article.ClientID(),
			"clientArticleId": article.ClientArticleId(),
			"title":           article.Title(),
			"content":         article.Content(),
			"lastModified":    article.LastModified(),
			"publishDate":     article.PublishDate(),
			"createdAt":       time.Now(), //For TTl
		}
		res, err := articlesCollection.InsertOne(ctx, repoArticle)
		if err != nil {
			return fmt.Errorf("error inserting new document: %v", err)
		}
		log.Info().Msg(fmt.Sprintf("Successfully inserted new article into database %v", res.InsertedID))
		return nil
	}

	if err != nil {
		return fmt.Errorf("error finding document: %v", err)
	}

	//If Article exists and lastmodifed is different we update the article
	if article.LastModified() > existingAritcle.LastModified {
		updateFields := bson.M{
			"title":        article.Title(),
			"content":      article.Content(),
			"lastModified": article.LastModified(),
			"publishDate":  article.PublishDate(),
		}
		res, err := articlesCollection.UpdateOne(ctx, bson.M{"_id": existingAritcle.Id}, updateFields)
		if err != nil {
			return fmt.Errorf("error updating document: %v", err)
		}
		log.Info().Msg(fmt.Sprintf("Successfully updated article in the database %v", res.UpsertedID))
	}

	return nil
}

// A function to get article by id
func (r ArticlesRepo) GetArticleById(ctx context.Context, id string, clientId string) (*domain.Article, error) {

	articlesCollection := r.client.Database(r.dbName).Collection(ArticlesCollection)

	var aritcle Article
	err := articlesCollection.FindOne(ctx, bson.M{"_id": id, "clientId": clientId}).Decode(&aritcle)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error finding document: %v", err) // other database error
	}

	//converting to domain article
	domainArticle, err := domain.NewArticle(
		aritcle.Id,
		aritcle.ClientId,
		aritcle.ClientArticleId,
		aritcle.Title,
		aritcle.Content,
		aritcle.LastModified,
		aritcle.PublishDate,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to covert to domain article %v", err)
	}

	return domainArticle, nil
}

// A function to get all articles
func (r ArticlesRepo) GetAllArticles(ctx context.Context, clientId string, limit int64, offset int64) ([]*domain.Article, error) {
	articlesCollection := r.client.Database(r.dbName).Collection(ArticlesCollection)

	// Create find options to limit results
	findOptions := options.Find()
	findOptions.SetLimit(limit) //setting limit to the number of articles to return
	findOptions.SetSkip(offset) // Setting offset to skip the specified number of articles

	// Finding documents
	cursor, err := articlesCollection.Find(ctx, bson.M{"clientId": clientId}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %v", err)
	}
	defer cursor.Close(ctx)

	var articles []*domain.Article
	for cursor.Next(ctx) {
		var article Article
		if err := cursor.Decode(&article); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}

		domainArticle, err := domain.NewArticle(
			article.Id,
			article.ClientId,
			article.ClientArticleId,
			article.Title,
			article.Content,
			article.LastModified,
			article.PublishDate,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to covert to domain article %v", err)
		}

		articles = append(articles, domainArticle)
	}

	// Check for cursor errors after the loop
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %v", err)
	}

	return articles, nil
}
