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

type RefreshTokensReqResp struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
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

	response := BaseResponse{Data: RefreshTokensReqResp{RefreshToken: string(tokens.RefreshToken), AccessToken: string(tokens.AccessToken)}, Message: ""}
	response_raw, _ := json.Marshal(response)
	writer.Write(response_raw)
	return
}

func (controller *AuthController) RefreshTokens(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var tokens RefreshTokensReqResp
	decoder := json.NewDecoder(reader.Body)
	if err := decoder.Decode(&tokens); err != nil {
		response := BaseResponse{Message: "Invalid JSON format"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}
	isValid := auth.ValidateTokensPair(auth.Token(tokens.AccessToken), auth.Token(tokens.RefreshToken))
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Invalid tokens pair"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	refreshTokenFromDb, get_token_err := controller.Repo.GetRefreshToken(tokens.RefreshToken)

	if get_token_err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error getting refresh token"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}
	if refreshTokenFromDb.Token != tokens.RefreshToken {
		writer.WriteHeader(http.StatusBadRequest)
		response := BaseResponse{Message: "Such token does not exist. Generate new pair"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	newAccessToken, generate_err := auth.GenerateAccessTokenByRefresh(auth.Token(tokens.RefreshToken))
	if generate_err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := BaseResponse{Message: "Error generating new access token"}
		response_raw, _ := json.Marshal(response)
		writer.Write(response_raw)
		return
	}

	response := BaseResponse{Data: RefreshTokensReqResp{AccessToken: string(newAccessToken), RefreshToken: tokens.RefreshToken}, Message: ""}
	response_raw, _ := json.Marshal(response)
	writer.Write(response_raw)
	return
}
