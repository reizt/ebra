package middlewares

import (
	"errors"
	"time"

	"github.com/reizt/ebra/conf"
	"github.com/reizt/ebra/models"
)

var (
	sessionDuration = time.Hour * 24 * 3
	// sessionDuration = time.Hour * 0
)

func GetCurrentUser(sessionId string) (currentUser *models.User, err error) {
	session := new(models.Session)
	resourcesDb := conf.ConnectMySQL()
	sessionDb := conf.ConnectSessionDB()
	err = sessionDb.First(&session, "id = ?", sessionId).Error
	if err != nil {
		return nil, err
	}
	expiresAt := session.CreatedAt.Add(sessionDuration)
	if time.Now().After(expiresAt) {
		return nil, errors.New("session expired")
	}

	err = resourcesDb.First(&currentUser, "id = ?", session.UserID).Error
	if currentUser == nil || err != nil {
		return nil, errors.New("current user not found")
	}
	return currentUser, nil
}

func StartSession(userId string) (sessionId string, err error) {
	db := conf.ConnectSessionDB()
	session := &models.Session{}
	session.UserID = userId
	// Delete existing session
	res := db.Where("user_id = ?", userId).Delete(&models.Session{})
	if res.Error != nil {
		return "", res.Error
	}
	res = db.Create(&session)
	if res.Error != nil {
		return "", res.Error
	}
	return session.ID, nil
}

func DisableSession(sessionId string) (err error) {
	db := conf.ConnectSessionDB()
	session := &models.Session{}
	res := db.Delete(&session, "id = ?", sessionId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
