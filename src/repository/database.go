package repository

import "github.com/jackc/pgx"

type RepoInterface interface {
	SetRefreshToken(refreshToken string) error
	GetRefreshToken(refreshToken string) (RefreshToken, error)
}

type DatabaseRepo struct {
	conn *pgx.Conn
}

func (database *DatabaseRepo) SetRefreshToken(refreshToken string) error {
	_, err := database.conn.Exec("INSERT INTO refresh_tokens (token) VALUES ($1)", refreshToken)
	return err
}

func (database *DatabaseRepo) GetRefreshToken(refreshToken string) (RefreshToken, error) {
	token_row, err := database.conn.Query("SELECT * FROM refresh_tokens WHERE token = $1", refreshToken)
	var token RefreshToken
	if err != nil {
		return token, err
	}

	for token_row.Next() {
		err = token_row.Scan(&token.Token, &token.UserID)
		if err != nil {
			return token, err
		}
	}
	return token, nil
}
