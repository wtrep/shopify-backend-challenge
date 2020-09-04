package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wtrep/image-repo-backend/api"
	"os"
	"sync"
)

type CreateImageError struct {
	image api.Image
	err   error
}

func NewConnectionPool() (*sql.DB, error) {
	dbIP, ok := os.LookupEnv("DB_IP")
	if !ok {
		return nil, errors.New("error: DB_IP environment variable not set")
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, errors.New("error: DB_PASSWORD environment variable not set")
	}

	db, err := sql.Open("mysql", "backend-sa:"+dbPassword+"@("+dbIP+")/images-repo")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateImage(db *sql.DB, image api.Image) error {
	uuid, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	stmt, err := generateCreateImageStatement(db)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(uuid, image.Name, image.Owner, image.Kind, image.Height, image.Length,
		image.Bucket, image.BucketPath, image.Status)
	if err != nil {
		return err
	}
	return nil
}

func CreateImages(db *sql.DB, images []api.Image) ([]CreateImageError, error) {
	stmt, err := generateCreateImageStatement(db)
	if err != nil {
		return nil, err
	}
	failedChan := make(chan CreateImageError, len(images))
	var wg sync.WaitGroup
	wg.Add(len(images))

	for _, image := range images {
		go createImageFromStmt(stmt, image, failedChan, &wg)
	}
	wg.Wait()
	close(failedChan)

	var createImageErrors []CreateImageError
	for e := range failedChan {
		createImageErrors = append(createImageErrors, e)
	}

	if len(createImageErrors) != 0 {
		return createImageErrors, errors.New("unable to create all images")
	}
	return nil, nil
}

func createImageFromStmt(stmt *sql.Stmt, image api.Image, failedChan chan<- CreateImageError, wg *sync.WaitGroup) {
	uuid, err := image.UUID.MarshalBinary()
	if err != nil {
		failedChan <- CreateImageError{
			image: image,
			err:   err,
		}
		wg.Done()
		return
	}

	_, err = stmt.Exec(uuid, image.Name, image.Owner, image.Kind, image.Height, image.Length,
		image.Bucket, image.BucketPath, image.Status)
	if err != nil {
		failedChan <- CreateImageError{
			image: image,
			err:   err,
		}
	}
	wg.Done()
}

func generateCreateImageStatement(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("INSERT INTO images " +
		"(UUID, name, owner, kind, height, length, bucket, bucketPath, status) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?, ?, ?, ?);")
}
