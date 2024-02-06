package controllers

import (
	"auth/src/auth"
	"auth/src/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type AuthController struct {
	Repo repository.RepoInterface
}

func (controller *AuthController) GenerateTokens(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	userIdFromRequest := reader.URL.Query().Get("user_id")
	if userIdFromRequest == "" {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Field user_id is required"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	userID, parse_err := uuid.Parse(userIdFromRequest)
	if parse_err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "User ID is not valid UUID"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	tokens, generate_err := auth.GenerateTokens(userID)
	if generate_err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error generating tokens"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	set_error := controller.Repo.SetRefreshToken(string(tokens.RefreshToken), userID)
	if set_error != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error setting refresh token"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	response := BaseResponse{Data: RefreshTokensReqResp{RefreshToken: string(tokens.RefreshToken), AccessToken: string(tokens.AccessToken)}, Message: ""}
	responseRaw, _ := json.Marshal(response)
	writer.Write(responseRaw)
	return
}

func (controller *AuthController) RefreshTokens(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var tokens RefreshTokensRequest
	decoder := json.NewDecoder(reader.Body)
	if err := decoder.Decode(&tokens); err != nil {
		response := BaseResponse{Message: "Invalid JSON format"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	isValid := auth.ValidateTokensPair(auth.Token(tokens.AccessToken), auth.Token(tokens.RefreshToken))
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Invalid tokens pair"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	tokenFromDb, get_token_err := controller.Repo.GetRefreshToken(tokens.RefreshToken)

	if get_token_err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Such token does not exist. Generate new pair"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	userId, _ := auth.GetUserIdFromToken(auth.Token(tokens.RefreshToken))
	if userId != tokenFromDb.UserID {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Invalid Refresh Token"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	newAccessToken, generate_err := auth.GenerateAccessTokenByRefresh(auth.Token(tokens.RefreshToken))
	if generate_err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error generating new access token"}
		responseRaw, _ := json.Marshal(response)
		writer.Write(responseRaw)
		return
	}

	response := BaseResponse{Data: RefreshTokensResponse{AccessToken: string(newAccessToken), RefreshToken: tokens.RefreshToken}, Message: ""}
	responseRaw, _ := json.Marshal(response)
	writer.Write(responseRaw)
	return
}
