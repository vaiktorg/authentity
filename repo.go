package authentity

import (
	"database/sql"
	"github.com/vaiktorg/Authentity/entities"
	"sync"

	"gorm.io/gorm"
)

type DBRepo struct {
	sync.Mutex
	db *gorm.DB
}

func NewAuthentityRepo(db *gorm.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// FindIdentityByID returns Identity when matched with a ProfileID.
func (a *DBRepo) FindIdentityByID(id string) (iden *entities.Identity, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&iden, "id = ?", id)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindIdentityByProfileID returns Identity when matched with a ProfileID.
func (a *DBRepo) FindIdentityByProfileID(profileId string) (iden *entities.Identity, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&iden, "profile_id = ?", profileId)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindIdentityByAccountID returns Identity when matched with a AccountID.
func (a *DBRepo) FindIdentityByAccountID(accId string) (iden *entities.Identity, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&iden, "account_id = ?", accId)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindAccountByUsername returns Account when matched with a username.
func (a *DBRepo) FindAccountByUsername(username string) (acc *entities.Account, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&acc, "username = ?", username)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindAccountByEmail returns Account when matched with a email.
func (a *DBRepo) FindAccountByEmail(email string) (acc *entities.Account, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&acc, "email = ?", email)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindProfileByUsername returns Profile when matched with a username.
func (a *DBRepo) FindProfileByUsername(username string) (acc *entities.Profile, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&acc, "username = ?", username)

		return nil
	}, &sql.TxOptions{})

	return
}

// FindProfileByEmail returns Profile when matched with a email.
func (a *DBRepo) FindProfileByEmail(email string) (acc *entities.Profile, err error) {
	a.Lock()
	defer a.Unlock()

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(&acc, "email = ?", email)

		return nil
	}, &sql.TxOptions{})

	return
}

// All returns array of result
func (a *DBRepo) All(dst []interface{}) (err error) {
	a.Lock()
	defer a.Unlock()

	return a.db.Transaction(func(tx *gorm.DB) error {
		tx.Find(dst)

		return nil
	}, &sql.TxOptions{})
}

// Find anything by ID in their gorm.Model struct.
func (a *DBRepo) Find(dst interface{}) (err error) {
	a.Lock()
	defer a.Unlock()

	return a.db.Transaction(func(tx *gorm.DB) error {
		tx.Find(dst)

		return nil
	}, &sql.TxOptions{})
}

// Persist saves if ID not found, and updates if ID found.
func (a *DBRepo) Persist(dst interface{}) error {
	a.Lock()
	defer a.Unlock()

	return a.db.Transaction(func(tx *gorm.DB) error {
		tx.Save(dst)
		return nil
	})
}
