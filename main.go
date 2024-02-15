package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	var port = os.Getenv("PORT")

	mainR := chi.NewRouter()
	mainR.Use(cors.Handler(corsOptions))

	v1R := chi.NewRouter()

	v1R.Get("/readiness", v1OkHandler)
	v1R.Get("/err", v1NotOkHandler)

	mainR.Mount("/v1", v1R)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mainR,
	}

	log.Printf("Starting http server on port %s", port)
	log.Fatal(srv.ListenAndServe())

}
