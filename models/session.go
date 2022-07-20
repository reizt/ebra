package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	Base
	UserID string
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New().String()
	return nil
}
