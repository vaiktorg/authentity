package entities

type Account struct {
	Model
	Username string // Unnique
	Email    string // Unnique
	Password string
}
