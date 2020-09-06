package common

import (
	"net/http"
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

var InvalidRequestBodyError = DetailedError{
	ID:     1200,
	Name:   "InvalidRequestBodyError",
	Detail: "The request body doesn't respect the valid format",
	Code:   http.StatusBadRequest,
}

var UserDoesNotExistError = DetailedError{
	ID:     1201,
	Name:   "UserDoesNotExistError",
	Detail: "The provided user doesn't exist",
	Code:   http.StatusNotFound,
}

var WrongPasswordError = DetailedError{
	ID:     1202,
	Name:   "WrongPasswordError",
	Detail: "The provided password doesn't match our records",
	Code:   http.StatusUnauthorized,
}

var DatabaseInsertionError = DetailedError{
	ID:     1203,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var JSONEncoderError = DetailedError{
	ID:     1204,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var PasswordTooLongError = DetailedError{
	ID:     1205,
	Name:   "PasswordTooLongError",
	Detail: "Password length is more than 32 characters",
	Code:   http.StatusBadRequest,
}

var UserAlreadyExistError = DetailedError{
	ID:     1206,
	Name:   "UserAlreadyExistError",
	Detail: "The provided user already exist",
	Code:   http.StatusConflict,
}

var MissingTokenError = DetailedError{
	ID:     1207,
	Name:   "MissingTokenError",
	Detail: "No bearer token was provided",
	Code:   http.StatusUnauthorized,
}

var InvalidTokenError = DetailedError{
	ID:     1208,
	Name:   "InvalidTokenError",
	Detail: "The bearer token provided is invalid",
	Code:   http.StatusBadRequest,
}

var TokenGenerationError = DetailedError{
	ID:     1209,
	Name:   "TokenGenerationError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusInternalServerError,
}

var TokenExpiredError = DetailedError{
	ID:     1210,
	Name:   "TokenExpiredError",
	Detail: "The provided token is expired",
	Code:   http.StatusUnauthorized,
}

var WrongUserError = DetailedError{
	ID:     1211,
	Name:   "WrongUserError",
	Detail: "The user associated with that token doesn't match the image owner",
	Code:   http.StatusUnauthorized,
}

var InvalidImageBodyError = DetailedError{
	ID:     1212,
	Name:   "InvalidImageBodyError",
	Detail: "Invalid image body. Maximum file size is 10 MB",
	Code:   http.StatusUnauthorized,
}

var FileUploadError = DetailedError{
	ID:     1213,
	Name:   "InternalServerError",
	Detail: "An unhandled error occurred, please try again",
	Code:   http.StatusUnauthorized,
}
