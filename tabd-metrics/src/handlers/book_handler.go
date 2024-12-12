package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/eduardor2m/tabd-metrics/src/models"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateBookHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// PostgreSQL

		db.Exec("INSERT INTO books (title, author, genre, description) VALUES ($1, $2, $3, $4)", book.Title, book.Author, book.Genre, book.Description)

		// MongoDB

		mongoDB.Collection("books").InsertOne(r.Context(), book)

		// Redis

		redisClient.Set(r.Context(), book.Title, book.Genre, 0)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Book created successfully"})
	}

}

func GetBooksHandler(db *sql.DB, mongoDB *mongo.Database, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT id, title, author, genre, description FROM books")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var books []models.Book
		for rows.Next() {
			var book models.Book
			if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Description); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			books = append(books, book)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}
