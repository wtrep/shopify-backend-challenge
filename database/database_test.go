package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/wtrep/image-repo-backend/api"
	"testing"
)

func TestCreateImage(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}

	image := api.Image{
		UUID:       uuid.New(),
		Name:       "testImage",
		Owner:      "William",
		Kind:       "jpeg",
		Height:     640,
		Length:     320,
		Bucket:     "testBucket",
		BucketPath: "william/testImage",
		Status:     "CREATED",
	}

	err = CreateImage(db, image)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDuplicateCreateImage(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}

	image := api.Image{
		UUID:       uuid.New(),
		Name:       "testImage",
		Owner:      "William",
		Kind:       "jpeg",
		Height:     640,
		Length:     320,
		Bucket:     "testBucket",
		BucketPath: "william/testImage",
		Status:     "CREATED",
	}

	err = CreateImage(db, image)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = CreateImage(db, image)
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number != 1062 {
			t.Errorf("Duplicate entry should return Error 1062")
		}
	}
}

func TestCreateImages(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}

	images := []api.Image{
		{
			UUID:       uuid.New(),
			Name:       "testImage",
			Owner:      "William",
			Kind:       "jpeg",
			Height:     640,
			Length:     320,
			Bucket:     "testBucket",
			BucketPath: "william/testImage",
			Status:     "CREATED",
		},
		{
			UUID:       uuid.New(),
			Name:       "testImage2",
			Owner:      "William",
			Kind:       "jpeg",
			Height:     6400,
			Length:     3200,
			Bucket:     "testBucket",
			BucketPath: "william/testImage",
			Status:     "CREATED",
		},
	}
	createImageErrors, err := CreateImages(db, images)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(createImageErrors) != 0 {
		t.Errorf("createImages return at least one createImageError")
	}
}

func TestDuplicateCreateImages(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}
	dummyUUID := uuid.New()

	images := []api.Image{
		{
			UUID:       dummyUUID,
			Name:       "testImage",
			Owner:      "William",
			Kind:       "jpeg",
			Height:     640,
			Length:     320,
			Bucket:     "testBucket",
			BucketPath: "william/testImage",
			Status:     "CREATED",
		},
		{
			UUID:       dummyUUID,
			Name:       "testImage2",
			Owner:      "William",
			Kind:       "jpeg",
			Height:     6400,
			Length:     3200,
			Bucket:     "testBucket",
			BucketPath: "william/testImage",
			Status:     "CREATED",
		},
	}
	createImageErrors, err := CreateImages(db, images)
	if err == nil {
		t.Errorf("duplicateCreateImages should return an error")
	}
	if len(createImageErrors) != 1 {
		t.Errorf("createImages should return one createImageError")
	}
}

func TestDeleteImage(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}

	image := api.Image{
		UUID:       uuid.New(),
		Name:       "testDeleteImage",
		Owner:      "William",
		Kind:       "jpeg",
		Height:     640,
		Length:     320,
		Bucket:     "testBucket",
		BucketPath: "william/testImage",
		Status:     "CREATED",
	}

	err = CreateImage(db, image)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = DeleteImage(db, image.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateImage(t *testing.T) {
	db, err := NewConnectionPool()
	if err != nil {
		t.Errorf(err.Error())
	}

	image1 := api.Image{
		UUID:       uuid.New(),
		Name:       "testUpdateImageFirst",
		Owner:      "William",
		Kind:       "jpeg",
		Height:     640,
		Length:     320,
		Bucket:     "testBucket",
		BucketPath: "william/testImage",
		Status:     "CREATED",
	}
	image2 := api.Image{
		UUID:       image1.UUID,
		Name:       "testUpdateImageSecond",
		Owner:      "Bill",
		Kind:       "png",
		Height:     6400,
		Length:     3200,
		Bucket:     "testBucketUpdated",
		BucketPath: "william/testImage123",
		Status:     "UPLOADED",
	}

	err = CreateImage(db, image1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = UpdateImage(db, image2)
	if err != nil {
		t.Errorf(err.Error())
	}
}
