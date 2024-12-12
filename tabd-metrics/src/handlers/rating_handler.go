package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/eduardor2m/tabd-metrics/src/models"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRatingHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rating models.Rating
		if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// PostgreSQL

		db.Exec("INSERT INTO ratings (user_id, book_id, note, comment) VALUES ($1, $2, $3, $4)", rating.UserID, rating.BookID, rating.Note, rating.Comment)

		// MongoDB

		mongoDB.Collection("ratings").InsertOne(r.Context(), rating)

		// Redis

		redisClient.Set(r.Context(), rating.UserID.String(), rating.BookID, 0)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Rating created successfully"})
	}

}

func GetRatingsHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT id, user_id, book_id, note, comment FROM ratings")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var ratings []models.Rating
		for rows.Next() {
			var rating models.Rating
			if err := rows.Scan(&rating.ID, &rating.UserID, &rating.BookID, &rating.Note, &rating.Comment); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ratings = append(ratings, rating)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ratings)
	}
}
