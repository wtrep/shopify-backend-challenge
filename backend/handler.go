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
	tokenLifetime = 2 * time.Hour
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
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandleGetAuthToken(w http.ResponseWriter, r *http.Request) {
	var sessionRequest UserRequest
	err := json.NewDecoder(r.Body).Decode(&sessionRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	currentUser, err := GetUser(h.db, sessionRequest.Username)
	if err != nil {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}

	if currentUser.goodPassword(sessionRequest.Password) {
		session := UserActiveSession{
			UUID:       uuid.New(),
			Username:   sessionRequest.Username,
			CreatedAt:  time.Now(),
			Expiration: time.Now().Add(tokenLifetime),
		}
		err := CreateActiveSession(h.db, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(session)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
	}
}

func (h *Handler) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(request.Password) > 32 {
		http.Error(w, "Password length is more than 32 characters", http.StatusBadRequest)
		return
	}
	user, err := NewUser(request.Username, request.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = CreateUser(h.db, *user)
	if err == UserAlreadyExist {
		http.Error(w, "User already exist", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func (h *Handler) HandlePostImage(w http.ResponseWriter, r *http.Request) {
	var request ImageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		http.Error(w, "Missing Bearer token", http.StatusBadRequest)
		return
	}
	tokenSlice := strings.Split(token, " ")
	if len(tokenSlice) != 2 {
		http.Error(w, "Invalid Bearer token format", http.StatusBadRequest)
		return
	}
	session, err := GetActiveSession(h.db, tokenSlice[1])
	if err != nil {
		http.Error(w, "No session associated with that token", http.StatusUnauthorized)
		return
	}

	image := request.toImage(session.Username)
	err = CreateImage(h.db, image)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(image.toImageResponse())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
