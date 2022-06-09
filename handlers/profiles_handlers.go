package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vaiktorg/authentity/entities"
	"net/http"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

// ProfilesHandler Return profiles
func ProfilesHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	bks, err := AllProfiles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(bks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id not provided", http.StatusInternalServerError)
		return
	}

	bks, err := FetchProfile(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(bks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//==================================

func FetchProfile(ctx context.Context, id string) (*entities.Profile, error) {
	// Retrieve the connection pool from the context. Because the
	// r.Context().Value() method always returns an interface{} type, we
	// need to type assert it into a *sql.DB before using it.
	db, ok := ctx.Value("access").(*DataProvider)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	prof := &entities.Profile{
		Model: entities.Model{ID: id},
	}

	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		tx.Find(prof)

		if prof == nil {
			return errors.New("could not fetch profile")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return prof, nil
}

func AllProfiles(ctx context.Context) ([]entities.Profile, error) {
	// Retrieve the connection pool from the context. Because the
	// r.Context().Value() method always returns an interface{} type, we
	// need to type assert it into a *sql.DB before using it.
	db, ok := ctx.Value("access").(*DataProvider)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	var bks []entities.Profile

	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		tx.Find(&bks)

		if len(bks) <= 0 {
			return errors.New("could not fetch profiles")
		}

		return nil
	}, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	return bks, nil
}
