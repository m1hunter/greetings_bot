package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //импорт драйвера
)

func Init() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func CreateTables(db *sql.DB) error {
	log.Println("Creating tables...")
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    username TEXT NOT NULL,
	    password TEXT NOT NULL,
	    first_name TEXT NOT NULL,
	    last_name TEXT NOT NULL,
	    birthday DATE NOT NULL,
	    chat_id BIGINT NOT NULL
	    );`

	createSubscriptionsTable := `
	CREATE TABLE IF NOT EXISTS subscriptions (
	    id SERIAL PRIMARY KEY,
	    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	    subscribed_to INT NOT NULL,
	    is_send_notification BOOLEAN NOT NULL
	    );`

	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("error creating users table: %v", err)
	}

	_, err = db.Exec(createSubscriptionsTable)
	if err != nil {
		log.Fatalf("error creating subscriprions table: %v", err)
	}

	log.Println("Tables created")
	return nil
}
