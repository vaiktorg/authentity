package authentity

import (
	"github.com/vaiktorg/authentity/entities"
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
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			&entities.Account{},
			&entities.Profile{},
			&entities.Identity{},
			&entities.Address{},
		)
	})
}
