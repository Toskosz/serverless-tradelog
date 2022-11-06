package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/Toskosz/everythingreviewed/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

func GetUserById(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	id := req.PathParameters["user-id"]
	userId, _ := strconv.Atoi(id)

	user, err := services.GetUserById(userId)
	if err != nil {
		return services.ApiResponse(http.StatusNotFound,
			api_error.NewNotFound("user", id))
	}

	return services.ApiResponse(http.StatusOK, user)
}

type loginInput struct {
	Email    string
	Password string
}

func (r *loginInput) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func Login(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	var jsonData loginInput

	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return services.ApiResponse(http.StatusBadRequest,
			api_error.NewBadRequest("Bad payload"))
	}

	jsonData.sanitize()
	user, err := services.Login(jsonData.Email, jsonData.Password)

	if err != nil {
		return services.ApiResponse(http.StatusBadRequest, err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("")))

	if err != nil {
		return services.ApiResponse(http.StatusInternalServerError,
			api_error.NewInternal())
	}

	cookie := []string{"jwt=" + token}

	return services.ApiResponseWithCookies(http.StatusOK, cookie, user)
}

type registerInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *registerInput) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func Register(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	var jsonData registerInput

	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return services.ApiResponse(http.StatusBadRequest,
			api_error.NewBadRequest("Bad payload"))
	}

	jsonData.sanitize()
	registerUserPayload := &models.User{
		Username: jsonData.Username,
		Email:    jsonData.Email,
		Password: jsonData.Password,
	}
	user, err := services.Register(registerUserPayload)

	if err != nil {
		if err.Error() == api_error.DuplicateEmailError {
			return services.ApiResponse(api_error.Status(err), err)
		} else {
			return services.ApiResponse(api_error.Status(err), err)
		}
	}

	return services.ApiResponse(http.StatusCreated, user)
}

func Logout(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	cookie := []string{"jwt=; Expires: " + time.Now().Add(-time.Hour).Format(
		"Mon, 02 Jan 2006 15:04:05 MST",
	)}
	return services.ApiResponseWithCookies(http.StatusOK, cookie, "Logged out")
}
