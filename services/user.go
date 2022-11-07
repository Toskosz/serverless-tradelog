package services

import (
	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
)

type userService struct {
	dbConn models.InterfaceDBUser
}

func NewUserService(conn models.InterfaceDBUser) models.InterfaceUserService {
	return &userService{
		dbConn: conn,
	}
}

func (s *userService) Login(email string, password string) (
	*models.User, error) {
	user, err := s.dbConn.FindUserByEmail(email)
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

func (s *userService) Register(user *models.User) (*models.User, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return nil, api_error.NewInternal()
	}

	user.Password = hash

	return s.dbConn.CreateUser(user)
}

func (s *userService) GetUserById(id int) (*models.User, error) {
	return s.dbConn.GetUserById(id)
}
