package backend

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	tokenLifetime = 2 * time.Hour
)

type Handler struct {
	db *sql.DB
}

func SetupRoutes(r *mux.Router) {
	db, err := NewConnectionPool()
	if err != nil {
		panic(err)
	}
	handler := Handler{db: db}

	r.HandleFunc("/auth", handler.HandleGetAuthToken).Methods("GET")
	//r.HandleFunc("/image/{uuid}", HandleGetImage).Methods("GET")
	http.Handle("/", r)
}

func (h Handler) HandleGetAuthToken(w http.ResponseWriter, r *http.Request) {
	var sessionRequest UserSessionRequest
	err := json.NewDecoder(r.Body).Decode(&sessionRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser, err := GetUser(h.db, sessionRequest.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentUser.goodPassword(sessionRequest.Password) {
		session := UserActiveSession{
			UUID:       uuid.New(),
			Username:   sessionRequest.Username,
			CreatedAt:  time.Time{},
			Expiration: time.Time{}.Add(tokenLifetime),
		}
		err := CreateActiveSession(h.db, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
