package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return nil
}
