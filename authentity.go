package main

import (
	"errors"
	"github.com/vaiktorg/Authentity/entities"
	"github.com/vaiktorg/Authentity/gwt"
	"strings"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type (
	Authentity struct {
		IdentityRepo *DBRepo
		Issuer       string
	}
)

var (
	Repo *Authentity
)

func RepoInit(dialector gorm.Dialector) {
	Repo = NewAuthentity(dialector)
}

func NewAuthentity(dialector gorm.Dialector) *Authentity {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	auth := &Authentity{
		IdentityRepo: NewAuthentityRepo(db),
		Issuer:       uuid.New().String(),
	}

	Migrate(db)

	return auth
}

func (a *Authentity) RegisterIdentity(prof *entities.Profile, acc *entities.Account) error {
	if _, err := a.IdentityRepo.FindAccountByEmail(acc.Email); err != nil {
		return errors.New("email already being used")
	}
	if _, err := a.IdentityRepo.FindAccountByUsername(acc.Username); err != nil {
		return errors.New("username already being used")
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not create hashed password")
	}

	identity := &entities.Identity{
		Profile: prof,
		Account: &entities.Account{
			Username: acc.Username,
			Email:    acc.Email,
			Password: string(hpass),
		},
	}

	err = a.IdentityRepo.Persist(identity)
	if err != nil {
		return err
	}

	return nil
}

func (a *Authentity) LoginToken(tkn string) error {
	gwt, err := gwt.DecodeGWT(tkn)
	if err != nil {
		return err
	}

	clearSig := func(iden *entities.Identity) error {
		iden.Signature = ""
		return a.IdentityRepo.Persist(iden)
	}

	if gwt.Header.Issuer != a.Issuer {
		return errors.New("token not issued by server")
	}

	if time.Since(gwt.Header.Timestamp) >= ExpireTime {
		return errors.New("token expired")
	}

	iden, err := a.IdentityRepo.FindIdentityByID(gwt.Header.ID)
	if err != nil {
		return errors.New("identity not found")
	}

	if iden.ID != gwt.Header.ID {
		return errors.New("identity id mismatch")
	}

	if strings.Compare(iden.Signature, gwt.Signature) != 0 {
		_ = clearSig(iden)
		return errors.New("signature mismatch")
	}

	return nil
}

func (a *Authentity) LoginManual(username, email, password string) (string, error) {
	acc, err := a.IdentityRepo.FindAccountByEmail(email)
	if err != nil {
		acc, err = a.IdentityRepo.FindAccountByUsername(username)
		if err != nil {
			return "", errors.New("models not found")
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		return "", errors.New("password does not match")
	}

	iden, err := a.IdentityRepo.FindIdentityByAccountID(acc.ID)
	if err != nil {
		return "", err
	}

	tok, sig, err := gwt.NewToken(iden.ID, a.Issuer, gwt.Spice.Salt, gwt.Spice.Pepper, iden)
	if err != nil {
		return "", err
	}

	iden.Signature = sig

	err = a.IdentityRepo.Persist(iden)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func (a *Authentity) LogoutToken(tkn string) error {
	gwt, err := gwt.DecodeGWT(tkn)
	if err != nil {
		return err
	}

	iden, err := a.IdentityRepo.FindIdentityByID(gwt.Header.ID)
	if err != nil {
		return errors.New("account not found")
	}

	iden.Signature = ""

	return a.IdentityRepo.Persist(iden)
}
