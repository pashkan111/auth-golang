package controllers

import (
	"auth/src/auth"
	"auth/src/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type BaseResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AuthController struct {
	Repo repository.RepoInterface
}

func (controller *AuthController) GenerateTokens(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	userIdFromRequest := reader.URL.Query().Get("user_id")
	if userIdFromRequest == "" {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Field user_id is required"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	userID, parse_err := uuid.Parse(userIdFromRequest)
	if parse_err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "User ID is not valid UUID"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	tokens, generate_err := auth.GenerateTokens(userID)
	if generate_err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error generating tokens"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	set_error := controller.Repo.SetRefreshToken(string(tokens.RefreshToken), userID)
	if set_error != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error setting refresh token"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	response := BaseResponse{Data: tokens}
	response_raw, _ := json.Marshal(response)
	writer.Write(response_raw)
	return
}
