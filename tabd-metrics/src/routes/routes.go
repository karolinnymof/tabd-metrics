package routes

import (
	"database/sql"

	"github.com/eduardor2m/tabd-metrics/src/handlers"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *mux.Router, db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) {
	r.HandleFunc("/users", handlers.CreateUserHandler(db, mongoDB, redisClient)).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsersHandler(db, mongoDB, redisClient)).Methods("GET")
	r.HandleFunc(("/users/performance"), handlers.GetUsersPerformanceHandler(db, mongoDB, redisClient)).Methods("GET")

	r.HandleFunc("/books", handlers.CreateBookHandler(db, mongoDB, redisClient)).Methods("POST")
	r.HandleFunc("/books", handlers.GetBooksHandler(db, mongoDB, redisClient)).Methods("GET")

	r.HandleFunc("/ratings", handlers.CreateRatingHandler(db, mongoDB, redisClient)).Methods("POST")
	r.HandleFunc("/ratings", handlers.GetRatingsHandler(db, mongoDB, redisClient)).Methods("GET")
}
