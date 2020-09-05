package repo_backend

import (
	"github.com/google/uuid"
)

type Image struct {
	UUID       uuid.UUID
	Name       string
	Owner      string
	Kind       string
	Height     int
	Length     int
	Bucket     string
	BucketPath string
	Status     string
}
