package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vaiktorg/authentity/entities"
	"net/http"

	"gorm.io/gorm"
)

// AccountsHandler Return profiles
func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	bks, err := AllAccounts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(bks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	// Pass the request context onto the database layer.
	info := &struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ID       string `json:"id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	bks, err := GetAccount(r.Context(), info.Username, info.Email)
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

func GetAccount(ctx context.Context, username, email string) (*entities.Account, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	var acc = &entities.Account{}
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Take(acc, "username = ? and email = ?", username, email)

		if acc == nil {
			return errors.New("could not fetch account")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func AllAccounts(ctx context.Context) ([]entities.Account, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not get database connection pool from context")
	}

	var bks []entities.Account
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Find(&bks)

		if len(bks) <= 0 {
			return errors.New("could not fetch accounts")
		}

		return nil
	}, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}

	return bks, nil
}
