package backend

import (
	"database/sql"
	"github.com/google/uuid"
)

type CreateImageError struct {
	image Image
	err   error
}

func CreateImage(db *sql.DB, image Image) error {
	uuidToCreate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO images (UUID, name, owner, kind, height, length, bucket, bucketPath, "+
		"status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuidToCreate, image.Name, image.Owner, image.Kind, image.Height,
		image.Length, image.Bucket, image.BucketPath, image.Status)
	if err != nil {
		return err
	}
	return nil
}

func CreateImages(db *sql.DB, images []Image) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, i := range images {
		uuidToCreate, err := i.UUID.MarshalBinary()
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("INSERT INTO images (UUID, name, owner, kind, height, length, bucket, bucketPath, "+
			"status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuidToCreate, i.Name, i.Owner, i.Kind, i.Height, i.Length,
			i.Bucket, i.BucketPath, i.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
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

func GetImage(db *sql.DB, id uuid.UUID) (*Image, error) {
	uuidToGet, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT * FROM images WHERE UUID = ?", uuidToGet)
	image := &Image{}
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

func UpdateImage(db *sql.DB, image Image) error {
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
