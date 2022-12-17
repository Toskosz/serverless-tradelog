package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/Toskosz/serverless-tradelog/services"
	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) GetUserByUsername(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	username := req.PathParameters["username"]

	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}

	return services.ApiResponse(http.StatusOK, user)
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *loginInput) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Username = strings.ToLower(r.Username)
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
	user, err := h.userService.Login(jsonData.Username, jsonData.Password)

	if err != nil {
		return services.ApiResponse(http.StatusBadRequest, err)
	}

	signedToken, err := h.userService.NewToken(user.Username)
	if err != nil {
		return services.ApiResponse(api_error.Status(err), err)
	}
	cookie := []string{"jwt=" + signedToken}

	return services.ApiResponseWithCookies(http.StatusOK, cookie, user)
}

type registerInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *registerInput) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Username = strings.ToLower(r.Username)
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
		Username:  jsonData.Username,
		Password:  jsonData.Password,
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}
	user, err := h.userService.Register(registerUserPayload)

	if err != nil {
		if err.Error() == api_error.DuplicateUsernameError {
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
		time.RFC3339Nano),
	}
	return services.ApiResponseWithCookies(http.StatusOK, cookie, "Logged out")
}
