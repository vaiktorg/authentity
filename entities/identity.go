package entities

import "github.com/google/uuid"

type Identity struct {
	Model

	ProfileID uuid.UUID // Unique
	Profile   *Profile  `gorm:"foreignKey:ProfileID"`

	AccountID string
	Account   *Account `gorm:"foreignKey:AccountID"`

	GroupsID string
	Groups   *Groups `gorm:"foreignKey:GroupsID"`

	PermissionsID string
	Permissions   *Permissions `gorm:"foreignKey:PermissionsID"`

	Signature string
}
