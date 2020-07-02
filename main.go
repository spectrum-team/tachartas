package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	gorillah "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spectrum-team/tachartas/commons"
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
	u := handlers.NewUserHandler(config)

	router := mux.NewRouter()

	// Events
	eventSubRouter(e, router)

	// Category
	router.HandleFunc("/category", c.FindAll).Methods("GET")

	// Auth
	router.HandleFunc("/signin", u.SignIn).Methods("POST")

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	listen := fmt.Sprintf(":%s", port)
	headersOk := gorillah.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "X-Requested-With", "User-Agent", "Referer", "Host", "Content-Type"})
	// originsOk := gorillah.AllowedOrigins([]string{"*"})
	methodsOk := gorillah.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(listen, gorillah.CORS(headersOk, methodsOk)(router))
}

func eventSubRouter(eventHandler *handlers.EventHandler, router *mux.Router) {

	eventRouter := router.PathPrefix("/event").Subrouter()
	eventRouter.Use(commons.AuthMiddleware)

	eventRouter.HandleFunc("/{id}", eventHandler.FindOne).Methods("GET")
	eventRouter.HandleFunc("/filter", eventHandler.Find).Methods("POST", "OPTIONS")
	eventRouter.HandleFunc("", eventHandler.Insert).Methods("POST", "OPTIONS")
	eventRouter.HandleFunc("/{id}", eventHandler.Update).Methods("PUT")
	eventRouter.HandleFunc("/{id}/image", eventHandler.AddImageToEvent).Methods("PUT")
	eventRouter.HandleFunc("/{id}/{assist}", eventHandler.Assist).Methods("PUT")
}
