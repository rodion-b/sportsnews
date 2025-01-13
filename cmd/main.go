package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sports-news-api/internal/app/config"
	"sports-news-api/internal/app/repository"
	"sports-news-api/internal/app/server"
	"sports-news-api/internal/app/services"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const refreshFrequency = time.Second * 10
const ecbArticlesPollSize = 5
const ecbFeedsRequestsTimeout = time.Second * 5

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("Error in run")
	}
	os.Exit(0)
}

func run() error {
	// read config from env
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create MongoDB client and connect
	clientOptions := options.Client().ApplyURI(cfg.MONGO_URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Err(err).Msg("Error connecting to MongoDB")
		return err
	}

	// Ping to ensure connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Err(err).Msg("Error pinging MongoDB")
		return err
	}
	log.Info().Msg("Successfully connected to DB")

	// Repositories Init
	articlesRepo := repository.NewArticlesRepo(client.Database(cfg.MONGO_DATABASE_NAME))
	// 	Services Init
	articlesService := services.NewArticlesService(articlesRepo)

	//Polling ECB Articles
	ecbFeeds := services.NewEcbFeedsService("content-ecb.pulselive.com", "content/ecb/text/EN/", ecbFeedsRequestsTimeout)
	go func() {
		ticker := time.NewTicker(refreshFrequency)
		for {
			select {
			case <-ticker.C:
				//get latest N articles
				ecbArticles, err := ecbFeeds.GetEcbArticles(ecbArticlesPollSize)
				if err != nil {
					log.Err(err).Msg("Unable to GetEcbArticles")
					continue //skipping the rest
				}
				for _, article := range ecbArticles.Content {
					go func() {
						//enriching article data
						articleFull, err := ecbFeeds.GetArticleById(strconv.Itoa(article.ID))
						if err != nil {
							log.Err(err).Msg("Unable to GetArticleById")
							return
						}
						//upserting article into db
						err = articlesService.UpsertEcbArticle(ctx, *articleFull)
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
	}()

	log.Info().Msg("Successfully started Polling ECB Articles")

	//Server Init
	server := server.NewServer(articlesService)

	// create http router
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Sports News API v0.1"))
	}).Methods("GET")

	//Setting up routes
	router.HandleFunc("/articles/{article_id}", server.GetArticleById).Methods(http.MethodGet)
	router.HandleFunc("/articles", server.GetAllArticles).Methods(http.MethodGet)

	//Http Server Init
	srv := &http.Server{
		Addr:    cfg.HTTP_ADDR,
		Handler: router,
	}

	// listen to OS signals and gracefully shutdown server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Err(err).Msg("HTTP Server Shutdown Error")
		}
		close(stopped)
	}()

	log.Info().Msg(fmt.Sprintf("Starting HTTP server on %s", cfg.HTTP_ADDR))

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped
	log.Info().Msg("Server stopped")

	return nil
}
