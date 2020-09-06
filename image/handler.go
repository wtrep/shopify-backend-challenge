package image

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wtrep/image-repo-backend/common"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	tokenLifetime = 24 * time.Hour
)

type Handler struct {
	db *sql.DB
}

func SetupAndServeRoutes() {
	db, err := NewConnectionPool()
	if err != nil {
		panic(err)
	}
	handler := Handler{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/image", handler.HandlePostImage).Methods("POST")
	r.HandleFunc("/image/{uuid}", handler.HandleGetImage).Methods("GET")
	//r.HandleFunc("/image/{uuid}", handler.HandleDeleteImage).Methods("DELETE")
	//r.HandleFunc("/images", handler.HandleGetImages).Methods("GET")
	r.HandleFunc("/upload/{uuid}", handler.HandlePostUpload).Methods("POST")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandlePostImage(w http.ResponseWriter, r *http.Request) {
	var request CreateImageRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	image := request.toImage(username)
	err = CreateImage(h.db, image)
	if err != nil {
		common.RespondWithError(w, &common.DatabaseInsertionError)
		return
	}

	err = json.NewEncoder(w).Encode(image.toCreateImageResponse())
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToGet, err := uuid.Parse(vars["uuid"])
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	image, err := GetImage(h.db, uuidToGet)
	if err != nil {
		common.RespondWithError(w, &common.ImageNotFoundError)
		return
	}

	if username == image.Owner && image.Status == "UPLOADED" {
		url, err := generateSignedURL(image.Bucket, image.BucketPath)
		if err != nil {
			common.RespondWithError(w, &common.URLGenerationError)
			return
		}

		response := image.toImageResponse(url)
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			common.RespondWithError(w, &common.JSONEncoderError)
		}
	} else if image.Status == "CREATED" {
		common.RespondWithError(w, &common.ImageNotUploadedError)
	} else {
		common.RespondWithError(w, &common.UserPermissionDeniedError)
	}
}

func (h *Handler) HandlePostUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToUpload, err := uuid.Parse(vars["uuid"])
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	image, err := GetImage(h.db, uuidToUpload)
	if err != nil {
		common.RespondWithError(w, &common.ImageNotFoundError)
		return
	}

	username, detailedErr := handleJWT(r)
	if detailedErr != nil {
		common.RespondWithError(w, detailedErr)
		return
	}
	if image.Owner != username {
		common.RespondWithError(w, &common.WrongUserError)
		return
	}

	file, detailedErr := getImageFromForm(w, r)
	if detailedErr != nil {
		common.RespondWithError(w, detailedErr)
		return
	}
	defer file.Close()

	err = uploadToBucket(file, image.Bucket, image.BucketPath)
	if err != nil {
		common.RespondWithError(w, &common.FileUploadError)
		return
	}

	image.Status = "UPLOADED"
	err = UpdateImage(h.db, *image)
	if err != nil {
		common.RespondWithError(w, &common.DatabaseInsertionError)
		return
	}

	url, err := generateSignedURL(image.Bucket, image.BucketPath)
	if err != nil {
		common.RespondWithError(w, &common.URLGenerationError)
	}

	response := image.toImageResponse(url)
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

func handleJWT(r *http.Request) (string, *common.ErrorResponseError) {
	if r.Header["Key"] == nil {
		return "", &common.MissingTokenError
	}

	username, err := common.VerifyJWT(r.Header["Key"][0])
	if err != nil {
		return "", &common.InvalidTokenError
	}
	return username, nil
}

func getImageFromForm(w http.ResponseWriter, r *http.Request) (multipart.File, *common.ErrorResponseError) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, &common.InvalidImageBodyError
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, &common.InvalidImageBodyError
	}
	return file, nil
}