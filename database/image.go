package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/wtrep/image-repo-backend/api"
	"sync"
)

type CreateImageError struct {
	image api.Image
	err   error
}

func CreateImage(db *sql.DB, image api.Image) error {
	uuidToCreate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	stmt, err := generateCreateImageStatement(db)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(uuidToCreate, image.Name, image.Owner, image.Kind, image.Height, image.Length,
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
	uuidToCreate, err := image.UUID.MarshalBinary()
	if err != nil {
		failedChan <- CreateImageError{
			image: image,
			err:   err,
		}
		wg.Done()
		return
	}

	_, err = stmt.Exec(uuidToCreate, image.Name, image.Owner, image.Kind, image.Height, image.Length,
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

func DeleteImage(db *sql.DB, id uuid.UUID) error {
	uuidToDelete, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM `images` WHERE `uuid` = ?", uuidToDelete)
	if err != nil {
		return err
	}
	return nil
}

func GetImage(db *sql.DB, id uuid.UUID) (*api.Image, error) {
	uuidToGet, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT * FROM images WHERE UUID = ?", uuidToGet)
	image := &api.Image{}
	var uuidToParse []byte

	err = row.Scan(&uuidToParse, &image.Name, &image.Owner, &image.Kind, &image.Height, &image.Length, &image.Bucket,
		&image.BucketPath, &image.Status)
	if err != nil {
		return nil, err
	}

	err = image.UUID.UnmarshalBinary(uuidToParse)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func UpdateImage(db *sql.DB, image api.Image) error {
	uuidToUpdate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE images SET name = ?, owner = ?, kind = ?, height = ?, length = ?, bucket = ?, "+
		"bucketPath = ?, status = ? WHERE uuid = ?", image.Name, image.Owner, image.Kind, image.Height, image.Length,
		image.Bucket, image.BucketPath, image.Status, uuidToUpdate)

	if err != nil {
		return err
	}
	return nil
}
