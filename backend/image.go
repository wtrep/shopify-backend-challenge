package backend

import (
	"github.com/google/uuid"
	"os"
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

type ImageRequest struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Height int    `json:"height"`
	Length int    `json:"length"`
}

type ImageResponse struct {
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Owner  string    `json:"owner"`
	Kind   string    `json:"kind"`
	Height int       `json:"height"`
	Length int       `json:"length"`
}

func (i *Image) toImageResponse() ImageResponse {
	return ImageResponse{
		UUID:   i.UUID,
		Name:   i.Name,
		Owner:  i.Owner,
		Kind:   i.Kind,
		Height: i.Height,
		Length: i.Length,
	}
}

func (i *ImageRequest) toImage(owner string) Image {
	imageUUID := uuid.New()

	return Image{
		UUID:       imageUUID,
		Name:       i.Name,
		Owner:      owner,
		Kind:       i.Kind,
		Height:     i.Height,
		Length:     i.Length,
		Bucket:     os.Getenv("BUCKET"),
		BucketPath: imageUUID.String() + "." + i.Kind,
		Status:     "CREATED",
	}
}
