package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agroconnect/service-catalog/handlers"
	"github.com/agroconnect/service-catalog/repository"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://agro_admin:agro_mongo_password@localhost:27017/agroconnect_catalog?authSource=admin"
	}

	mongoDBName := os.Getenv("MONGO_DB")
	if mongoDBName == "" {
		mongoDBName = "agroconnect_catalog"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Warning: MongoDB ping failed: %v", err)
	} else {
		log.Println("Successfully connected to MongoDB")
	}

	db := client.Database(mongoDBName)
	productRepo := repository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service":  "agroconnect-service-catalog",
			"status":   "UP",
			"database": "MongoDB",
			"version":  "1.0.0",
		})
	}).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")
	api.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	api.HandleFunc("/products/{id}", productHandler.GetProductByID).Methods("GET")
	api.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
	api.HandleFunc("/products/{id}/stock", productHandler.DeductProductStock).Methods("PATCH")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)

	log.Printf("AgroConnect Service Catalog running on port :%s", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server terminated: %v", err)
	}
}
