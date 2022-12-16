package services

import (
	"os"

	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/golang-jwt/jwt"
)

type userService struct {
	dbConn models.InterfaceDBUser
}

func NewUserService(conn models.InterfaceDBUser) models.InterfaceUserService {
	return &userService{
		dbConn: conn,
	}
}

func (s *userService) Login(username string, password string) (
	*models.User, error) {
	user, err := s.dbConn.FetchUserByUsername(username)
	if err != nil {
		return nil, api_error.NewAuthorization(api_error.InvalidCredentialsError)
	}

	match, err := comparePasswords(user.Password, password)
	if err != nil {
		return nil, api_error.NewUnsupportedMediaType(err.Error())
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

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.dbConn.FetchUserByUsername(username)
}

func (s *userService) GetUserFromToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", api_error.NewAuthorization("failed to get user from token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["Issuer"].(string)
		return username, nil
	} else {
		return "", api_error.NewAuthorization("failed to get user from token")
	}
}
