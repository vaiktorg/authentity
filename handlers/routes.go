package handlers

import (
	"context"
	"net/http"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func Routes(r *mux.Router) error {
	db, err := gorm.Open(sqlite.Open("DB.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	r.HandleFunc("/identities", injectDB(db, IdentitiesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/identity/{id}", injectDB(db, IdentityHandler)).Methods(http.MethodGet)

	r.HandleFunc("/profiles", injectDB(db, ProfilesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/profile/{id}", injectDB(db, ProfileHandler)).Methods(http.MethodGet)

	r.HandleFunc("/accounts", injectDB(db, AccountsHandler)).Methods(http.MethodGet)
	r.HandleFunc("/account{id}", injectDB(db, AccountHandler)).Methods(http.MethodPost)

	return nil
}

func injectDB(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
