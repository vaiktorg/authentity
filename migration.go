package main

import (
	"encoding/json"
	"github.com/vaiktorg/Authentity/entities"
	"os"
	"path/filepath"

	"github.com/brianvoe/gofakeit"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(
			&entities.Identity{},
			&entities.UserActivityLog{})
	})
}

func Drop(db *gorm.DB) error {
	tables := []interface{}{
		&entities.Account{},
		&entities.Profile{},
		&entities.Identity{},
		&entities.Groups{},
		&entities.Permissions{},
		&entities.Address{},
	}
	for _, model := range tables {
		err := db.Migrator().DropTable(model)
		if err != nil {
			return err
		}
	}
	return nil
}

func DummyData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var idens []entities.Account

		for i := 0; i < 20; i++ {
			username, email := gofakeit.Username(), gofakeit.Email()
			strPwd := gofakeit.Password(true, true, true, true, false, 8)
			hashPwd, err := bcrypt.GenerateFromPassword([]byte(strPwd), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			idens = append(idens, entities.Account{
				Username: username,
				Email:    email,
				Password: strPwd,
			})

			addr := gofakeit.Address()

			tx.Create(&entities.Identity{
				Profile: &entities.Profile{
					FirstName:   gofakeit.FirstName(),
					Initial:     gofakeit.Letter(),
					LastName:    gofakeit.LastName(),
					LastName2:   gofakeit.LastName(),
					PhoneNumber: gofakeit.Phone(),
					Address: &entities.Address{
						Addr1:   addr.Address,
						Addr2:   "n/a",
						City:    addr.City,
						State:   addr.State,
						Country: addr.Country,
						Zip1:    addr.Zip,
					},
				},
				Account: &entities.Account{
					Username: username,
					Email:    email,
					Password: string(hashPwd),
				},
			})
		}

		file, err := os.Create(filepath.Join(".", "LoginAccessFile.json"))
		if err != nil {
			return err
		}

		enc := json.NewEncoder(file)
		enc.SetIndent("", "\t")
		return enc.Encode(idens)
	})
}
