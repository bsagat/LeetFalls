package dbrepo

import (
	"database/sql"
	"errors"
	"leetFalls/internal/domain/models"
)

type AuthRepo struct {
	Db *sql.DB
}

func (repo *AuthRepo) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	err := repo.Db.QueryRow(`
		SELECT ID 
			FROM Users 
			WHERE Token_ID = $1 LIMIT 1;
	`, sessionID).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

// Changes user's name if it modified
func (repo *AuthRepo) ChangeUserName(id int, changedName string) error {
	_, err := repo.Db.Exec(`
	UPDATE 
		Users
	SET 
		Name = $1
	WHERE 
		id=$2 AND Name != $1;
	`, changedName, id)
	if err != nil {
		return err
	}
	return nil
}

// SaveUser saves a new user to the database
func (repo *AuthRepo) SaveUser(user models.User) error {
	_, err := repo.Db.Exec(`
	INSERT INTO
		Users (ID, Name, Token_ID, Avatar_URL)
	VALUES
		($1, $2, $3, $4)
	`, user.ID, user.Name, user.Token_ID, user.ImageURL)
	if err != nil {
		return err
	}
	return nil
}

// Gets unique user id
func (repo *AuthRepo) GetNextUserId() (int, error) {
	var id int
	if err := repo.Db.QueryRow("SELECT COALESCE(MAX(ID), 0) + 1 FROM Users").Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Gets user by id
func (repo *AuthRepo) GetUserById(id int) (models.User, error) {
	var user models.User
	err := repo.Db.QueryRow(`SELECT ID, name, Token_ID, TokenDate, Expires_at, Avatar_URL 
	FROM Users
	WHERE ID=$1
	`, id).Scan(&user.ID, &user.Name, &user.Token_ID, &user.TokenDate, &user.Expires_at, &user.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
