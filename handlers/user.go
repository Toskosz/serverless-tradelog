package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/Toskosz/serverless-tradelog/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) GetUserById(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	id := req.PathParameters["user-id"]
	userId, _ := strconv.Atoi(id)

	user, err := h.userService.GetUserById(userId)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
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

func (h *Handler) Login(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	var jsonData loginInput

	if err := json.Unmarshal([]byte(req.Body), &jsonData); err != nil {
		return services.ApiResponse(http.StatusBadRequest,
			api_error.NewBadRequest("Bad payload"))
	}

	jsonData.sanitize()
	user, err := h.userService.Login(jsonData.Email, jsonData.Password)

	if err != nil {
		return services.ApiResponse(http.StatusBadRequest, err)
	}

	// new jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))

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

func (h *Handler) Register(req events.APIGatewayProxyRequest) (
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
	user, err := h.userService.Register(registerUserPayload)

	if err != nil {
		if err.Error() == api_error.DuplicateEmailError {
			return services.ApiResponse(api_error.Status(err), err)
		} else {
			return services.ApiResponse(api_error.Status(err), err)
		}
	}

	return services.ApiResponse(http.StatusCreated, user)
}

func (h *Handler) Logout(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	cookie := []string{"jwt=; Expires: " + time.Now().Add(-time.Hour).Format(
		"Mon, 02 Jan 2006 15:04:05 MST",
	)}
	return services.ApiResponseWithCookies(http.StatusOK, cookie, "Logged out")
}
