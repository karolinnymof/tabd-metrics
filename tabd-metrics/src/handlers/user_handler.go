package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eduardor2m/tabd-metrics/src/models"
	"github.com/eduardor2m/tabd-metrics/src/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateUserHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// PostgreSQL

		db.Exec("INSERT INTO users (name, email, preferences) VALUES ($1, $2, $3)", user.Name, user.Email, user.Preferences)

		// MongoDB

		mongoDB.Collection("users").InsertOne(r.Context(), user)

		// Redis

		redisClient.Set(r.Context(), user.Email, user.Preferences, 0)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	}
}

func GetUsersHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT id, name, email, preferences FROM users")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Preferences); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		// MongoDB

		cursor, err := mongoDB.Collection("users").Find(r.Context(), bson.D{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(r.Context())

		for cursor.Next(r.Context()) {
			var user models.User
			if err := cursor.Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		// Redis

		keys := redisClient.Keys(r.Context(), "*").Val()
		for _, key := range keys {
			value := redisClient.Get(r.Context(), key).Val()
			user := models.User{Email: key, Preferences: value}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUsersPerformanceHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTotal := time.Now()

		startGenerate := time.Now()
		fakeUsers, err := utils.GenerateFakeUsers(1000)
		generateDuration := time.Since(startGenerate)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		startPostgres := time.Now()
		for _, user := range fakeUsers {
			_, err := db.Exec("INSERT INTO users (name, email, preferences) VALUES ($1, $2, $3)", user.Name, user.Email, user.Preferences)
			if err != nil {
				http.Error(w, fmt.Sprintf("PostgreSQL error: %v", err), http.StatusInternalServerError)
				return
			}
		}
		postgresDuration := time.Since(startPostgres)

		startMongo := time.Now()
		for _, user := range fakeUsers {
			_, err := mongoDB.Collection("users").InsertOne(r.Context(), user)
			if err != nil {
				http.Error(w, fmt.Sprintf("MongoDB error: %v", err), http.StatusInternalServerError)
				return
			}
		}
		mongoDuration := time.Since(startMongo)

		startRedis := time.Now()
		for _, user := range fakeUsers {
			err := redisClient.Set(r.Context(), user.Email, user.Preferences, 0).Err()
			if err != nil {
				http.Error(w, fmt.Sprintf("Redis error: %v", err), http.StatusInternalServerError)
				return
			}
		}
		redisDuration := time.Since(startRedis)

		totalDuration := time.Since(startTotal)

		report := map[string]string{
			"generate_users": fmt.Sprintf("%v", generateDuration),
			"postgres_time":  fmt.Sprintf("%v", postgresDuration),
			"mongo_time":     fmt.Sprintf("%v", mongoDuration),
			"redis_time":     fmt.Sprintf("%v", redisDuration),
			"total_time":     fmt.Sprintf("%v", totalDuration),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Users created successfully",
			"report":  report,
		})
	}
}
