package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Token string

type Claims struct {
	UserID uuid.UUID
	jwt.StandardClaims
	TokenAssociation uuid.UUID
}

type UserTokens struct {
	AccessToken  Token
	RefreshToken Token
}

func GenerateAccessToken(userID, tokenAssociation uuid.UUID) (Token, error) {
	expirationTime := time.Now().Add(time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		TokenAssociation: tokenAssociation,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return Token(signedToken), nil
}

func GenerateAccessTokenByRefresh(refreshToken Token) (Token, error) {
	refreshClaims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}
	return GenerateAccessToken(refreshClaims.UserID, refreshClaims.TokenAssociation)
}

func GenerateRefreshToken(userID, tokenAssociation uuid.UUID) (Token, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		TokenAssociation: tokenAssociation,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return Token(signedToken), nil
}

func ValidateToken(tokenString Token) (*Claims, error) {
	token, err := jwt.ParseWithClaims(string(tokenString), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

func ValidateTokensPair(accessToken, refreshToken Token) bool {
	accessClaims, err := ValidateToken(accessToken)
	if err != nil {
		return false
	}

	refreshClaims, err := ValidateToken(refreshToken)
	if err != nil {
		return false
	}

	return accessClaims.TokenAssociation == refreshClaims.TokenAssociation
}

func GenerateTokens(userID uuid.UUID) (*UserTokens, error) {
	tokenAssociation := uuid.New()
	accessToken, access_token_err := GenerateAccessToken(userID, tokenAssociation)
	if access_token_err != nil {
		return nil, access_token_err
	}
	refreshToken, refresh_token_err := GenerateRefreshToken(userID, tokenAssociation)
	if refresh_token_err != nil {
		return nil, refresh_token_err
	}
	return &UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GetUserIdFromToken(token Token) (uuid.UUID, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}
