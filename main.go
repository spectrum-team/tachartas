package main

import (
	"context"
	"fmt"
	gorillah "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spectrum-team/tachartas/handlers"
	"github.com/spectrum-team/tachartas/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
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

	router := mux.NewRouter()

	router.HandleFunc("/event/{id}", e.FindOne).Methods("GET")
	router.HandleFunc("/event/filter", e.Find).Methods("POST")
	router.HandleFunc("/event", e.Insert).Methods("POST")
	router.HandleFunc("/event/{id}", e.Update).Methods("PUT")

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	listen := fmt.Sprintf(":%s", port)

	http.ListenAndServe(listen, gorillah.CombinedLoggingHandler(os.Stdout, router))
}
