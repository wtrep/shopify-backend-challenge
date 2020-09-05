package backend

import (
	"github.com/google/uuid"
	"os"
)

type ErrorResponse struct {
	Error DetailedError `json:"error"`
}

type DetailedError struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
	Code   int    `json:"-"`
}

type SuccessfulCreateUserResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
}

type SuccessfulGenericResponse struct {
	Status string `json:"status"`
}

type CreateImageRequest struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Height int    `json:"height"`
	Length int    `json:"length"`
}

type SuccessfulCreateImageResponse struct {
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Owner  string    `json:"owner"`
	Kind   string    `json:"kind"`
	Height int       `json:"height"`
	Length int       `json:"length"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetTokenRequest = CreateUserRequest

func (i *Image) toImageResponse() SuccessfulCreateImageResponse {
	return SuccessfulCreateImageResponse{
		UUID:   i.UUID,
		Name:   i.Name,
		Owner:  i.Owner,
		Kind:   i.Kind,
		Height: i.Height,
		Length: i.Length,
	}
}

func (i *CreateImageRequest) toImage(owner string) Image {
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
