package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `gorm:"-:all"`
	PasswordDigest string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	bcryptCost := 12
	digest, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(digest)
	return nil
}
