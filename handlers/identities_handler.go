package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/vaiktorg/authentity/entities"
	"net/http"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

// IdentitiesHandler Return identities
func IdentitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	bks, err := FetchIdentities(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(bks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//IdentityHandler ...
func IdentityHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id not provided", http.StatusInternalServerError)
		return
	}

	bks, err := FetchIdentity(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(bks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func FetchIdentities(ctx context.Context) ([]entities.Identity, error) {
	access, ok := ctx.Value("access").(*DataProvider)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	var prof []entities.Identity

	access.Mutex.Lock()
	defer access.Mutex.Unlock()
	err := access.DB.Transaction(func(tx *gorm.DB) error {
		tx.Find(prof)

		if len(prof) <= 0 {
			return errors.New("could not fetch identities")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return prof, nil
}
func FetchIdentity(ctx context.Context, ID string) ([]entities.Identity, error) {
	access, ok := ctx.Value("db").(*DataProvider)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	var prof []entities.Identity

	access.Mutex.Lock()
	defer access.Mutex.Unlock()
	err := access.DB.Transaction(func(tx *gorm.DB) error {
		tx.Take(prof, "id = ?", ID)

		if prof == nil {
			return errors.New("could not fetch identity")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return prof, nil
}
