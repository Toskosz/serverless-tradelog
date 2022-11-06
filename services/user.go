package services

import (
	"github.com/Toskosz/everythingreviewed/db"
	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
)

func Login(email string, password string) (*models.User, error) {
	user, err := db.FindUserByEmail(email)
	if err != nil {
		return nil, api_error.NewAuthorization(api_error.InvalidCredentialsError)
	}

	match, err := comparePasswords(user.Password, password)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	if !match {
		return nil, api_error.NewAuthorization(api_error.InvalidCredentialsError)
	}

	return user, nil
}

func Register(user *models.User) (*models.User, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	user.Password = hash

	return db.CreateUser(user)
}

func GetUserById(id int) (*models.User, error) {
	return db.GetUserById(id)
}
