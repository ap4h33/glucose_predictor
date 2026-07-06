package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	//also add the refs to db
)

func main() {

	godotenv.Load("../.env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the env")
	}

	//Connect db with the env and shi: dbURL, conn(connection)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	//here goes the links
	v1Router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test works"))
	})
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Print("The server is active")
	err := srv.ListenAndServe() // When adding db change := to =
	if err != nil {
		log.Fatal(err)
	}
}
