package image

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"os"
)

func NewConnectionPool() (*sql.DB, error) {
	dbIP := os.Getenv("DB_IP")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.Getenv("DB_USERNAME")

	db, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@("+dbIP+")/"+dbName+
		"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateImage(db *sql.DB, image Image) error {
	uuidToCreate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO images (UUID, name, owner, extension, height, length, bucket, bucketPath, "+
		"status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuidToCreate, image.Name, image.Owner, image.Extension, image.Height,
		image.Length, image.Bucket, image.BucketPath, image.Status)
	if err != nil {
		return err
	}
	return nil
}

func UpdateImage(db *sql.DB, image Image) error {
	uuidToUpdate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE images SET name = ?, owner = ?, extension = ?, height = ?, length = ?, bucket = ?, "+
		"bucketPath = ?, status = ? WHERE uuid = ?", image.Name, image.Owner, image.Extension, image.Height, image.Length,
		image.Bucket, image.BucketPath, image.Status, uuidToUpdate)

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

	err = row.Scan(&uuidToParse, &image.Name, &image.Owner, &image.Extension, &image.Height, &image.Length, &image.Bucket,
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
