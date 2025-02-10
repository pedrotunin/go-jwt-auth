package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedrotunin/go-jwt-auth/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	log.Print("starting api...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("connecting to database")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("connected to database")

	router := gin.Default()
	router.Use(gin.Recovery())

	app := &config.Application{
		Router: router,
		DB:     db,
	}
	app.Setup()

	port := os.Getenv("PORT")

	log.Printf("running api on port %s", port)
	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("error executing server: %v", err)
	}
}
