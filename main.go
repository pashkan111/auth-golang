package main

import (
	"auth/src/controllers"
	"auth/src/repository"
	"log"
	"net/http"
)

func main() {
	repo := repository.DatabaseRepo{Conn: Db()}
	controller := controllers.AuthController{Repo: &repo}
	http.HandleFunc("/generate-tokens", controller.GenerateTokens)
	http.HandleFunc("/refresh-tokens", controller.RefreshTokens)
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
