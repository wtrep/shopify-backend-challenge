package auth

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wtrep/image-repo-backend/common"
	"net/http"
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
	r.HandleFunc("/user", handler.HandlePostUser).Methods("POST")
	r.HandleFunc("/key", handler.HandleGetKey).Methods("GET")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, &common.InvalidRequestBodyError)
		return
	}

	if len(request.Password) > 32 {
		respondWithError(w, &common.PasswordTooLongError)
		return
	}
	user, err := NewUser(request.Name, request.Password)
	if err != nil {
		respondWithError(w, &common.DatabaseInsertionError)
		return
	}

	err = CreateUser(h.db, *user)
	if err == UserAlreadyExist {
		respondWithError(w, &common.UserAlreadyExistError)
		return
	} else if err != nil {
		respondWithError(w, &common.DatabaseInsertionError)
		return
	}

	sendJWTToken(w, user)
}

func (h *Handler) HandleGetKey(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, &common.InvalidRequestBodyError)
		return
	}

	currentUser, err := GetUser(h.db, request.Name)
	if err != nil {
		respondWithError(w, &common.UserDoesNotExistError)
		return
	}

	if currentUser.goodPassword(request.Password) {
		sendJWTToken(w, currentUser)
	} else {
		respondWithError(w, &common.WrongPasswordError)
	}
}

func sendJWTToken(w http.ResponseWriter, user *User) {
	token, err := common.GenerateJWT(user.Username)
	if err != nil {
		respondWithError(w, &common.TokenGenerationError)
	}
	response := UserResponse{Token: token}

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		respondWithError(w, &common.JSONEncoderError)
	}
}

func respondWithError(w http.ResponseWriter, error *common.DetailedError) {
	w.WriteHeader(error.Code)
	response := common.ErrorResponse{
		Error: *error,
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		http.Error(w, "Server unhandled error", error.Code)
	}
}
