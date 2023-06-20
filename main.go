package main

import (
	"log"
	"net/http"
	"os"
	"userService/db/postgres"
	"userService/handler"
	"userService/user/service"
	"userService/user/storage"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	config := postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		Sslmode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := postgres.NewPostgresConnection(config)
	if err != nil {
		log.Fatalln(err)
	}

	userStorage := storage.NewUserStorage(db)
	authController := service.NewAuthService(userStorage)

	h := handler.NewHanadler(authController)

	http.ListenAndServe("localhost:8080", h.NewApiRouter())
}
