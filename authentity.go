package authentity

import (
	"errors"
	"github.com/vaiktorg/authentity/entities"
	"github.com/vaiktorg/gwt"
	"strings"
	"time"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type (
	Authentity struct {
		IdentityRepo *DBRepo
		Issuer       string
		encoder      *gwt.Encoder
		decoder      *gwt.Decoder
		db           *gorm.DB
	}

	AuthMan struct {
		login chan entities.Account
	}
)

func NewAuthentity(issuerName string, dialector gorm.Dialector) *Authentity {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	auth := &Authentity{
		IdentityRepo: NewAuthentityRepo(db),
		Issuer:       issuerName,
		db:           db,
		decoder:      gwt.NewDecoder(),
		encoder:      gwt.NewEncoder(),
	}

	if err = auth.Migrate(); err != nil {
		switch err {
		case AlreadyExistError:
			break
		}
	}

	return auth
}

func (a *Authentity) RegisterIdentity(prof *entities.Profile, acc *entities.Account) error {
	if _, err := a.IdentityRepo.FindAccountByEmail(acc.Email); err != nil {
		return errors.New("email already being used")
	}
	if _, err := a.IdentityRepo.FindAccountByUsername(acc.Username); err != nil {
		return errors.New("username already being used")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not create hashed password")
	}

	identity := &entities.Identity{
		Profile: prof,
		Account: &entities.Account{
			Username: acc.Username,
			Email:    acc.Email,
			Password: string(hashedPassword),
		},
	}

	err = a.IdentityRepo.Persist(identity)
	if err != nil {
		return err
	}

	return nil
}

func (a *Authentity) LoginToken(tkn gwt.Token) error {
	var errs error

	a.decoder.Decode(tkn, func(value gwt.Value, err error) {
		if err != nil {
			errs = err
			return
		}

		clearSig := func(identity *entities.Identity) error {
			identity.Signature = ""
			return a.IdentityRepo.Persist(identity)
		}

		if value.Issuer != a.Issuer {
			errs = errors.New("token not issued by server")
			return
		}

		if time.Since(value.Timestamp) >= ExpireTime {
			errs = errors.New("token expired")
			return
		}

		account, err := a.IdentityRepo.FindAccountByUsername(value.Username)
		if err != nil {
			errs = errors.New("identity not found")
		}

		identity, err := a.IdentityRepo.FindIdentityByAccountID(account.ID)
		if err != nil {
			errs = errors.New("identity not found")
		}

		if strings.Compare(identity.Signature, string(value.Signature)) != 0 {
			_ = clearSig(identity)
			errs = errors.New("signature mismatch")
		}
	})

	return errs
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

	identity, err := a.IdentityRepo.FindIdentityByAccountID(acc.ID)
	if err != nil {
		return "", err
	}

	var tok gwt.Token
	a.encoder.Encode(gwt.Value{
		Issuer:    a.Issuer,
		Username:  acc.Username,
		Timestamp: time.Now().Add(time.Hour),
	}, func(token gwt.Token, err error) {
		identity.Signature = string(token.Signature)
		e := a.IdentityRepo.Persist(identity)
		if e != nil {
			err = e
		}

		tok = token
	})
	if err != nil {
		return "", err
	}

	return tok.Token, nil
}

func (a *Authentity) LogoutToken(tkn string) error {
	var errs error
	a.decoder.Decode(gwt.Token{Token: tkn}, func(value gwt.Value, err error) {
		if err != nil {
			errs = err
		}

		account, e := a.IdentityRepo.FindAccountByUsername(value.Username)
		if e != nil {
			errs = errors.New("account not found")
		}

		identity, e := a.IdentityRepo.FindIdentityByAccountID(account.ID)
		if e != nil {
			errs = errors.New("account not found")
		}

		identity.Signature = ""

		e = a.IdentityRepo.Persist(identity)
		if err != nil {
			errs = err
		}
	})

	return errs
}
