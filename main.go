package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jabuta/feedreader/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	var conf apiConfig

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	var port = os.Getenv("PORT")
	var dbURL = os.Getenv("DB_CONNECTION_STRING")

	if db, err := sql.Open("postgres", dbURL); err != nil {
		panic(err)
	} else {
		conf.DB = database.New(db)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()

	//Define bolerplate endpoints
	v1router.Get("/readiness", v1OkHandler)
	v1router.Get("/err", v1NotOkHandler)

	//Define Users endpoints
	v1router.Post("/users", conf.createUser)
	v1router.Get("/users", conf.middlewareAuth(conf.handlerUsersGet))

	// define feeds endpoints
	v1router.Post("/feeds", conf.middlewareAuth(conf.handlerFeedsPost))
	v1router.Get("/feeds", conf.handlerFeedsGet)

	// define feed_follow endpoints
	v1router.Post("/feed_follows", conf.middlewareAuth(conf.createFeedFollow))
	v1router.Get("/feed_follows", conf.middlewareAuth(conf.getFeedFollowsUser))
	v1router.Delete("/feed_follows/{feedFollowID}", conf.deleteFeedFollow)

	router.Mount("/v1", v1router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting http server on port %s", port)
	log.Fatal(srv.ListenAndServe())

}
