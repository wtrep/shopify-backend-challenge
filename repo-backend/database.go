package repo_backend

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func NewConnectionPool() (*sql.DB, error) {
	dbIP, ok := os.LookupEnv("DB_IP")
	if !ok {
		return nil, errors.New("error: DB_IP environment variable not set")
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, errors.New("error: DB_PASSWORD environment variable not set")
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("error: DB_NAME environment variable not set")
	}

	db, err := sql.Open("mysql", "repo-backend-sa:"+dbPassword+"@("+dbIP+")/"+dbName+
		"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}
