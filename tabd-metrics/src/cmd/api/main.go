package main

import (
	"log"
	"net/http"

	"github.com/eduardor2m/tabd-metrics/src/db"
	"github.com/eduardor2m/tabd-metrics/src/routes"
	"github.com/gorilla/mux"
)

func main() {
	postgresDB, err := db.InitPostgres()
	if err != nil {
		log.Fatal("Erro ao conectar ao PostgreSQL:", err)
	}
	defer postgresDB.Close()

	mongoDB, err := db.InitMongoDB()
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}

	redisClient := db.InitRedis()
	defer redisClient.Close()

	r := mux.NewRouter()
	routes.RegisterRoutes(r, postgresDB, mongoDB, redisClient)

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// package main

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand/v2"
// 	"sync"
// 	"time"

// 	"github.com/go-faker/faker/v4"
// 	"github.com/go-redis/redis/v8"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	_ "github.com/lib/pq"
// )

// var idCounter int
// var idMutex sync.Mutex

// type User struct {
// 	ID          int    `json:"id"`
// 	Name        string `faker:"username"`
// 	Email       string `faker:"email"`
// 	Preferences string `faker:"word"`
// }

// type Book struct {
// 	ID          int    `json:"id"`
// 	Title       string `faker:"word"`
// 	Author      string `faker:"word"`
// 	Genre       string `faker:"word"`
// 	Description string `faker:"sentence"`
// }

// type Rating struct {
// 	ID      int    `json:"id"`
// 	UserID  int    `json:"user_id"`
// 	BookID  int    `json:"book_id"`
// 	Note    int    `json:"note"`
// 	Comment string `faker:"sentence"`
// }

// type Data struct {
// 	Users   []User
// 	Books   []Book
// 	Ratings []Rating
// }

// func createDataFake(numUsers, numBooks, numRatings int) ([]Data, error) {

// 	users := make([]User, numUsers)
// 	if err := faker.FakeData(&users); err != nil {
// 		return nil, fmt.Errorf("failed to generate users: %w", err)
// 	}

// 	books := make([]Book, numBooks)
// 	if err := faker.FakeData(&books); err != nil {
// 		return nil, fmt.Errorf("failed to generate books: %w", err)
// 	}

// 	ratings := make([]Rating, numRatings)
// 	if err := faker.FakeData(&ratings); err != nil {
// 		return nil, fmt.Errorf("failed to generate ratings: %w", err)
// 	}

// 	data := []Data{
// 		{Users: users, Books: books, Ratings: ratings},
// 	}

// 	return data, nil
// }

// // insertDataToPostgres inserts user data into a PostgreSQL database.
// func insertDataToPostgres(db *sql.DB, users []User) error {
// 	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, name VARCHAR(50), email VARCHAR(50), preferences TEXT)")
// 	if err != nil {
// 		return fmt.Errorf("failed to create table: %w", err)
// 	}

// 	for _, user := range users {
// 		_, err := db.Exec("INSERT INTO users (name, email, preferences) VALUES ($1, $2, $3)", user.Name, user.Email, user.Preferences)
// 		if err != nil {
// 			return fmt.Errorf("failed to insert user data into Postgres: %w", err)
// 		}
// 	}
// 	return nil
// }

// // insertDataToMongoDB inserts user data into a MongoDB database.
// func insertDataToMongoDB(collection *mongo.Collection, users []User) error {
// 	for _, user := range users {
// 		_, err := collection.InsertOne(context.TODO(), user)
// 		if err != nil {
// 			return fmt.Errorf("failed to insert user data into MongoDB: %w", err)
// 		}
// 	}
// 	return nil
// }

// // insertDataToRedis inserts user data into a Redis database.
// func insertDataToRedis(rdb *redis.Client, users []User) error {
// 	for _, user := range users {
// 		userData, err := json.Marshal(user)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal user data to JSON: %w", err)
// 		}
// 		err = rdb.Set(context.TODO(), fmt.Sprintf("user:%d", user.ID), userData, 0).Err()
// 		if err != nil {
// 			return fmt.Errorf("failed to insert user data into Redis: %w", err)
// 		}
// 	}
// 	return nil
// }

// func main() {
// 	rand.Shuffle(time.Now().Nanosecond(), func(i, j int) {
// 		idMutex.Lock()
// 		idCounter++
// 		idMutex.Unlock()
// 	})
// 	data, err := createDataFake(5, 10, 15) // Generate 5 users, 10 books, and 15 ratings
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(data)
// 	usersData := []User{
// 		{1, "Alice", "alice@gmail.com", `{"theme": "dark", "language": "en"}`},
// 		{2, "Bob", "bob@gmail.com", `{"theme": "light", "language": "pt"}`},
// 	}

// 	// Postgres connection
// 	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/root?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// MongoDB connection
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/mydb?authSource=admin"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	collection := client.Database("mydb").Collection("mycollection")

// 	// Redis connection
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})
// 	defer rdb.Close()

// 	start := time.Now()

// 	// Insert data into all databases
// 	if err := insertDataToPostgres(db, usersData); err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := insertDataToMongoDB(collection, usersData); err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := insertDataToRedis(rdb, usersData); err != nil {
// 		log.Fatal(err)
// 	}

// 	end := time.Now()
// 	elapsed := end.Sub(start)
// 	fmt.Printf("Tempo de inserção: %s\n", elapsed)
// }
