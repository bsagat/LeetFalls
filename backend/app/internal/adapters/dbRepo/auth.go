package dbrepo

import (
	"database/sql"
	"errors"
)

type SessionRepo struct {
	Db *sql.DB
}

func (repo *SessionRepo) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	err := repo.Db.QueryRow(`
		SELECT ID FROM Users WHERE Token_ID = $1 LIMIT 1;
	`, sessionID).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}
