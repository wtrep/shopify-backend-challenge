package backend

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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
	r.HandleFunc("/auth", handler.HandleGetAuthToken).Methods("GET")
	r.HandleFunc("/user", handler.HandlePostUser).Methods("POST")
	r.HandleFunc("/image", handler.HandlePostImage).Methods("POST")
	r.HandleFunc("/images", handler.HandlePostImages).Methods("POST")
	r.HandleFunc("/upload/{uuid}", handler.HandlePostUpload).Methods("POST")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandleGetAuthToken(w http.ResponseWriter, r *http.Request) {
	var userRequest GetTokenRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		respondWithError(w, &InvalidRequestBodyError)
		return
	}

	currentUser, err := GetUser(h.db, userRequest.Username)
	if err != nil {
		respondWithError(w, &UserDoesNotExistError)
		return
	}

	if currentUser.goodPassword(userRequest.Password) {
		session, detailedErr := h.createUserSession(userRequest)
		if detailedErr != nil {
			respondWithError(w, detailedErr)
			return
		}

		err := json.NewEncoder(w).Encode(session)
		if err != nil {
			respondWithError(w, &JSONEncoderError)
		}
	} else {
		respondWithError(w, &WrongPasswordError)
	}
}

func (h *Handler) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, &InvalidRequestBodyError)
		return
	}

	if len(request.Password) > 32 {
		respondWithError(w, &PasswordTooLongError)
		return
	}
	user, err := NewUser(request.Username, request.Password)
	if err != nil {
		respondWithError(w, &DatabaseInsertionError)
		return
	}

	err = CreateUser(h.db, *user)
	if err == UserAlreadyExist {
		respondWithError(w, &UserAlreadyExistError)
		return
	} else if err != nil {
		respondWithError(w, &DatabaseInsertionError)
		return
	}

	response := SuccessfulCreateUserResponse{
		Status:   "OK",
		Username: user.Username,
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		respondWithError(w, &JSONEncoderError)
	}
}

func (h *Handler) HandlePostImage(w http.ResponseWriter, r *http.Request) {
	var request CreateImageRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, &InvalidRequestBodyError)
		return
	}

	session, detailedError := h.verifySessionToken(r)
	if detailedError != nil {
		respondWithError(w, detailedError)
		return
	}

	image := request.toImage(session.Username)
	err = CreateImage(h.db, image)
	if err != nil {
		respondWithError(w, &DatabaseInsertionError)
		return
	}

	err = json.NewEncoder(w).Encode(image.toImageResponse())
	if err != nil {
		respondWithError(w, &JSONEncoderError)
	}
}

func (h *Handler) HandlePostImages(w http.ResponseWriter, r *http.Request) {
	var request CreateImagesRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, &InvalidRequestBodyError)
		return
	}

	session, detailedError := h.verifySessionToken(r)
	if detailedError != nil {
		respondWithError(w, detailedError)
		return
	}

	images := request.toImages(session.Username)
	err = CreateImages(h.db, images)
	if err != nil {
		respondWithError(w, &DatabaseInsertionError)
		return
	}

	err = json.NewEncoder(w).Encode(generateImagesResponse(images))
	if err != nil {
		respondWithError(w, &JSONEncoderError)
	}
}

func (h *Handler) HandlePostUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToUpload, err := uuid.Parse(vars["uuid"])
	if err != nil {
		http.Error(w, "Invalid Image UUID Format", http.StatusBadRequest)
		return
	}

	image, err := GetImage(h.db, uuidToUpload)
	if err != nil {
		http.Error(w, "No image was found", http.StatusNotFound)
		return
	}

	session, detailedErr := h.verifySessionToken(r)
	if detailedErr != nil {
		respondWithError(w, detailedErr)
		return
	}
	if image.Owner != session.Username {
		respondWithError(w, &WrongUserError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, &InvalidImageBodyError)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, &InvalidImageBodyError)
		return
	}
	defer file.Close()

	err = UploadToBucket(file, image.Bucket, image.BucketPath)
	if err != nil {
		respondWithError(w, &FileUploadError)
		return
	}

	response := SuccessfulGenericResponse{
		Status: "OK",
	}
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		respondWithError(w, &JSONEncoderError)
	}
}

func (h *Handler) createUserSession(userRequest CreateUserRequest) (*UserActiveSession, *DetailedError) {
	session := UserActiveSession{
		UUID:       uuid.New(),
		Username:   userRequest.Username,
		CreatedAt:  time.Now(),
		Expiration: time.Now().Add(tokenLifetime),
	}
	err := CreateActiveSession(h.db, session)
	if err != nil {
		return nil, &DatabaseInsertionError
	}
	return &session, nil
}

func (h *Handler) verifySessionToken(r *http.Request) (*UserActiveSession, *DetailedError) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return nil, &MissingTokenError
	}
	tokenSlice := strings.Split(token, " ")
	if len(tokenSlice) != 2 {
		return nil, &InvalidTokenError
	}
	session, err := GetActiveSession(h.db, tokenSlice[1])
	if err != nil {
		return nil, &NoSessionError
	}
	if time.Now().After(session.Expiration) {
		return nil, &TokenExpiredError
	}
	return session, nil
}

func respondWithError(w http.ResponseWriter, error *DetailedError) {
	w.WriteHeader(error.Code)
	response := ErrorResponse{
		Error: *error,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		http.Error(w, "Server unhandled error", error.Code)
	}
}
