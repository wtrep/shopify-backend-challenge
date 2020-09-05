package backend

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var UserAlreadyExist = errors.New("user already exist")

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec("INSERT INTO users (username, pwHash) VALUEs (?, ?)", user.Username, user.PwHash)
	if err, ok := err.(*mysql.MySQLError); ok {
		if err.Number == 1062 {
			return UserAlreadyExist
		}
		return errors.New(err.Error())
	}
	return nil
}

func UpdateUser(db *sql.DB, user User) error {
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

func GetUser(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	user := &User{}

	err := row.Scan(&user.Username, &user.PwHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateActiveSession(db *sql.DB, session UserActiveSession) error {
	uuidToCreate, err := session.UUID.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO sessions (UUID, username, created, expiration) VALUES (?, ?, ?, ?)",
		uuidToCreate, session.Username, session.CreatedAt, session.Expiration)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSessionExpiration(db *sql.DB, session UserActiveSession) error {
	uuidToFind, err := session.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE sessions SET expiration = ? WHERE UUID = ?", session.Expiration, uuidToFind)
	if err != nil {
		return err
	}
	return nil
}

func GetActiveSession(db *sql.DB, sessionUUID string) (*UserActiveSession, error) {
	uuidToFind, err := uuid.Parse(sessionUUID)
	if err != nil {
		return nil, err
	}
	binaryUUID, err := uuidToFind.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var session UserActiveSession
	row := db.QueryRow("SELECT * FROM sessions WHERE UUID=?", binaryUUID)
	err = row.Scan(&session.UUID, &session.Username, &session.CreatedAt, &session.Expiration)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
