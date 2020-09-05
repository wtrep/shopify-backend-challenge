package backend

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
