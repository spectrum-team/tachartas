package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	gorillah "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spectrum-team/tachartas/handlers"
	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClient(conn string) (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db := client.Database("tachartas")
	return db, nil
}

func main() {
	fmt.Println("Tachartas")

	conn := os.Getenv("DB_CONN_STRING")

	if conn == "" {
		conn = "mongodb://localhost:27017"
	}

	db, err := getMongoClient(conn)
	if err != nil {
		panic(err)
	}

	config := &models.DatabaseConfig{
		MongoClient: db,
	}

	e := handlers.NewEventHandler(config)
	c := handlers.NewCategoryHandler(config)

	router := mux.NewRouter()

	router.HandleFunc("/event/{id}", e.FindOne).Methods("GET")
	router.HandleFunc("/event/filter", e.Find).Methods("POST", "OPTIONS")
	router.HandleFunc("/event", e.Insert).Methods("POST", "OPTIONS")
	router.HandleFunc("/event/{id}", e.Update).Methods("PUT")
	router.HandleFunc("/event/{id}/image", e.AddImageToEvent).Methods("PUT")
	router.HandleFunc("/event/{id}/{assist}", e.Assist).Methods("PUT")

	// Category
	router.HandleFunc("/category", c.FindAll).Methods("GET")

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	listen := fmt.Sprintf(":%s", port)
	headersOk := gorillah.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "X-Requested-With"})
	// originsOk := gorillah.AllowedOrigins([]string{"*"})
	methodsOk := gorillah.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(listen, gorillah.CombinedLoggingHandler(os.Stdout, gorillah.CORS(headersOk, methodsOk)(router)))
}
