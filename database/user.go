package database

import (
	"database/sql"
	"github.com/wtrep/image-repo-backend/api"
)

func CreateUser(db *sql.DB, user api.User) error {
	_, err := db.Exec("INSERT INTO users (username, pwHash) VALUEs (?, ?)", user.Username, user.PwHash)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *sql.DB, user api.User) error {
	_, err := db.Exec("UPDATE users SET pwHash = ? WHERE username = ?", user.PwHash, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, username string) error {
	_, err := db.Exec("DELETE FROM `users` WHERE `username` = ?", username)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(db *sql.DB, username string) (*api.User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	user := &api.User{}

	err := row.Scan(&user.Username, &user.PwHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
