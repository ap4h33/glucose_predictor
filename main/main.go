package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //for sql
)

type apiConfig struct {
	DB         *database.Queries // the db package will be created by sqlc automatically
	HTTPClient *http.Client
	ModelURL   string // connection with the model
	AIVersion  string
	AImodelID  uuid.UUID
	ODUVersion string
	ODUmodelID uuid.UUID
	RecURL     string
}

func main() {

	godotenv.Load("../.env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in the env")
	}

	modelURL := os.Getenv("MODEL_URL")
	if modelURL == "" {
		log.Fatal("MODEL_URL not found in the env")
	}

	aiVersion := os.Getenv("AI_MODEL_VERSION")
	if aiVersion == "" {
		log.Fatal("AI_MODEL_VERSION not found in the env")
	}

	aiModelIDStr := os.Getenv("AI_MODEL_ID")
	if aiModelIDStr == "" {
		log.Fatal("AI_MODEL_ID not found in the env")
	}
	aiModelID, err := uuid.Parse(aiModelIDStr)
	if err != nil {
		log.Fatalf("invalid AI_MODEL_ID: %v", err)
	}

	oduVersion := os.Getenv("ODU_MODEL_VERSION")
	if oduVersion == "" {
		log.Fatal("ODU_MODEL_VERSION not found in the env")
	}

	oduModelIDStr := os.Getenv("ODU_MODEL_ID")
	if oduModelIDStr == "" {
		log.Fatal("ODU_MODEL_ID not found in the env")
	}
	oduModelID, err := uuid.Parse(oduModelIDStr)
	if err != nil {
		log.Fatalf("invalid ODU_MODEL_ID: %v", err)
	}

	recURL := os.Getenv("REC_URL")
	if recURL == "" {
		log.Fatal("REC_URL not found in the env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connection to database")
	}

	apiCfg := apiConfig{
		DB:         database.New(conn),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ModelURL:   modelURL,
		AIVersion:  aiVersion,
		AImodelID:  aiModelID,
		ODUVersion: oduVersion,
		ODUmodelID: oduModelID,
		RecURL:     recURL,
	} //this is used for hooking up links

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://", "http://localhost:3000*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/hospitals", apiCfg.handlerCreateHospital)
	v1Router.Get("/hospitals", apiCfg.handlerGetHospitals)
	v1Router.Get("/hospitals/{name}", apiCfg.handlerGetHospitals)

	v1Router.Post("/admin", apiCfg.handlerCreateAdmin)
	v1Router.Get("/admin/{admin_id}", apiCfg.handlerGetAdmin)

	v1Router.Get("/admin/user/{admin_id}", apiCfg.handlerGetUsers)

	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user/{patient_id}", apiCfg.handlerGetUser)

	v1Router.Post("/readings", apiCfg.handlerAddReadings)
	v1Router.Get("/glucose_levels/{patient_id}", apiCfg.handlerSeeInfo) // Full patient history
	v1Router.Get("/glucose_levels/{patient_id}/{time_period}", apiCfg.handlerSeeInfo)
	v1Router.Post("/predictions", apiCfg.handlerAddPredictions)
	v1Router.Get("/predictions", apiCfg.handlerGetAllModelPredictions)
	//GET latest predictions are in the glucose_levels, along with patient history

	v1Router.Get("/recommendations/{patient_id}", apiCfg.handlerGetRecommendations)
	v1Router.Post("/recommendations", apiCfg.handlerAddRecommendations)

	v1Router.Post("/models", apiCfg.handlerAddModel)
	// TO DO: handler get all models

	v1Router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test works"))
	})
	router.Mount("/glucose_predictor/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Print("The server is active")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
