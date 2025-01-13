package repository

import (
	"context"
	"fmt"
	"sports-news-api/internal/app/models"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticlesRepo struct {
	db *mongo.Database
}

func NewArticlesRepo(db *mongo.Database) *ArticlesRepo {
	return &ArticlesRepo{db: db}
}

func (r ArticlesRepo) UpsertEcbArticle(ctx context.Context, ecbArticle models.EcbArticle) error {
	filter := bson.M{
		"id": ecbArticle.ID, // assuming `ID` is a field in GetEcbArticleResponse
	}

	update := bson.M{
		"$set": ecbArticle, // updates or inserts with the entire struct
	}

	// Perform the upsert operation
	_, err := r.db.Collection(articles).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to upsert article: %v", err)
	}

	log.Info().Msg(fmt.Sprintf("Successfully upserted article with ID: %d", ecbArticle.ID))

	return nil
}

func (r ArticlesRepo) GetEcbArticleById(ctx context.Context, id string) (interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var result bson.M
	err = r.db.Collection(articles).FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error finding document: %v", err)
	}

	return result, nil
}

func (r ArticlesRepo) GetAllEcbArticles(ctx context.Context, limit int64) ([]interface{}, error) {
	// Create find options to limit results and sort by lastModified descending
	findOptions := options.Find()
	findOptions.SetLimit(limit) //setting limit to the number of articles to return

	// Finding documents
	cursor, err := r.db.Collection(articles).Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %v", err)
	}
	defer cursor.Close(ctx)

	var results []interface{}
	for cursor.Next(ctx) {
		var article bson.M
		if err := cursor.Decode(&article); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		results = append(results, article)
	}

	// Check for cursor errors after the loop
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %v", err)
	}

	return results, nil
}
