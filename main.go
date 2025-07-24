package main

import (
	"database/sql"
	"fmt"
	"golang-final-project/config"
	"golang-final-project/database"
	"golang-final-project/router"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	startServer()
}

func connectToDB() {
	config.Err = godotenv.Load("config/.env")
	if config.Err != nil {
		panic("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		// os.Getenv("PGHOST"),
		// os.Getenv("PGPORT"),
		// os.Getenv("PGUSER"),
		// os.Getenv("PGPASSWORD"),
		// os.Getenv("PGDATABASE"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	config.Db, config.Err = sql.Open("postgres", psqlInfo)
	if config.Err != nil {
		panic(config.Err)
	}
	// defer Db.Close()

	config.Err = config.Db.Ping()
	if config.Err != nil {
		panic(config.Err)
	}
	database.DBMigrate(config.Db)

	fmt.Println("Successfully connected to the database")
}

func startServer() {
	var PORT = "8080"
	// var PORT = os.Getenv("PORT")

	connectToDB()
	router.StartServer().Run(":" + PORT)
}
