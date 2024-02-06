package controllers

type BaseResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type RefreshTokensResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
