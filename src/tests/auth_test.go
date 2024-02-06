package tests

import (
	"auth/src/auth"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateTokensValidated(t *testing.T) {
	userId := uuid.New()
	tokens, _ := auth.GenerateTokens(userId)

	validated := auth.ValidateTokensPair(tokens.AccessToken, tokens.RefreshToken)
	if validated != true {
		t.Errorf("Tokens are not valid")
	}
}

func TestGenerateTokensValidationError(t *testing.T) {
	userId := uuid.New()
	tokenAssociation := uuid.New()

	tokens, _ := auth.GenerateTokens(userId)
	otherAccessToken, _ := auth.GenerateAccessToken(userId, tokenAssociation)

	validated := auth.ValidateTokensPair(otherAccessToken, tokens.RefreshToken)
	if validated != false {
		t.Errorf("Tokens are not valid")
	}
}
