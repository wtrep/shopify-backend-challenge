package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/wtrep/image-repo-backend/api"
	"os"
	"reflect"
	"testing"
)

var (
	backupDatabase string
	db             *sql.DB
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	returnCode := m.Run()
	tearDown()
	os.Exit(returnCode)
}

func setup() error {
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return errors.New("error: DB_NAME environment variable not set")
	}
	backupDatabase = dbName
	err := os.Setenv("DB_NAME", "images-repo-test")
	if err != nil {
		return err
	}

	db, err = NewConnectionPool()
	return err
}

func tearDown() {
	_, err := db.Exec("DELETE FROM images")
	if err != nil {
		fmt.Println("warning: unable to clear table images")
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		fmt.Println("warning: unable to clear table users")
	}
	err = os.Setenv("DB_NAME", backupDatabase)
	if err != nil {
		fmt.Println("warning: unable to restore DB_NAME environment variable")
	}
}

func TestCreateImage(t *testing.T) {
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

	err := CreateImage(db, image)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDuplicateCreateImage(t *testing.T) {
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

	err := CreateImage(db, image)
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

	err := CreateImage(db, image)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = DeleteImage(db, image.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUpdateImage(t *testing.T) {
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

	err := CreateImage(db, image1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = UpdateImage(db, image2)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetImage(t *testing.T) {
	image1 := &api.Image{
		UUID:       uuid.New(),
		Name:       "testGetImage",
		Owner:      "William",
		Kind:       "jpeg",
		Height:     640,
		Length:     320,
		Bucket:     "testBucket",
		BucketPath: "william/testImage",
		Status:     "CREATED",
	}

	err := CreateImage(db, *image1)
	if err != nil {
		t.Errorf(err.Error())
	}

	image2, err := GetImage(db, image1.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(image2, image1) {
		t.Errorf("the two images should be identicals")
	}
}

func TestCreateUser(t *testing.T) {
	user, err := api.NewUser("dummyUser", "dummyPassword")
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("user shouldn't be nil")
	}

	err = CreateUser(db, *user)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDestroyUser(t *testing.T) {
	user, err := api.NewUser("dummyUserToDelete", "dummyPassword")
	if err != nil {
		t.Errorf(err.Error())
	}
	if user == nil {
		t.Errorf("user shouldn't be nil")
	}

	err = CreateUser(db, *user)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = DeleteUser(db, user.Username)
	if err != nil {
		t.Errorf(err.Error())
	}
}
