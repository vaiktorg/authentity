package tests

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/vaiktorg/authentity"
	"github.com/vaiktorg/authentity/entities"
	"github.com/vaiktorg/gwt"
	"os"
	"testing"
)

var (
	DBName      = "test_db"
	DbDialector = sqlite.Open(DBName)
	Auth        = authentity.NewAuthentity("TestAuthentity", DbDialector)
	TestProfile = &entities.Profile{
		FirstName:   "John",
		Initial:     "E",
		LastName:    "Smith",
		LastName2:   "Johnson",
		PhoneNumber: "1234567890",
		Address: &entities.Address{
			Addr1:   "666 Hellsing Ave.",
			Addr2:   "Bldg. 3x6 Apt. 543",
			City:    "Gehenna",
			State:   "3rd Circle",
			Country: "Inferno",
			Zip:     "00666",
		},
	}

	TestAccount = &entities.Account{
		Username: "nyarlathotep",
		Email:    "space-worm@elder1s.com",
		Password: "MrN00dle$123",
	}

	Token = gwt.Token{}
)

func TestMain(m *testing.M) {
	if Auth == nil {
		panic("auth is nil")
	}

	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		panic("sql db not created")
	}

	m.Run()
}

func TestRegister(t *testing.T) {
	tkn, err := Auth.LoginManual(
		TestAccount.Username,
		TestAccount.Email,
		TestAccount.Password,
	)
	if err != nil {

		err := Auth.RegisterIdentity(
			TestProfile,
			TestAccount,
		)
		if err != nil {
			t.FailNow()
		}
	}

	Token.Token = tkn
	fmt.Println(Token)
}

func TestLogin(t *testing.T) {
	if Token.Token == "" {
		tkn, err := Auth.LoginManual(
			TestAccount.Username,
			TestAccount.Email,
			TestAccount.Password,
		)
		if err != nil {
			t.FailNow()
		}

		Token.Token = tkn
		fmt.Println(Token)
	}
}

func TestLoginToken(t *testing.T) {
	err := Auth.LoginToken(Token)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}

func TestLogoutToken(t *testing.T) {
	err := Auth.LogoutToken(Token.Token)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}
