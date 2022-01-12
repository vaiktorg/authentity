package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vaiktorg/grimoire/helpers"

	"github.com/vaiktorg/Authentity/handlers"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func Routes(r *mux.Router) error {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf(`Authentity/%s.db`, helpers.MakeTimestampStr())), &gorm.Config{})
	if err != nil {
		return err
	}

	r.HandleFunc("/identities", injectDB(db, handlers.IdentitiesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/identity/{id}", injectDB(db, handlers.IdentityHandler)).Methods(http.MethodGet)

	r.HandleFunc("/profiles", injectDB(db, handlers.ProfilesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/profile/{id}", injectDB(db, handlers.ProfileHandler)).Methods(http.MethodGet)

	r.HandleFunc("/accounts", injectDB(db, handlers.AccountsHandler)).Methods(http.MethodGet)
	r.HandleFunc("/account{id}", injectDB(db, handlers.AccountHandler)).Methods(http.MethodPost)

	return nil
}

func injectDB(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
