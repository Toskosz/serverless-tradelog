package handlers

import (
	"net/http"
	"regexp"

	"github.com/Toskosz/serverless-tradelog/db"
	"github.com/Toskosz/serverless-tradelog/models"
	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/Toskosz/serverless-tradelog/services"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	userService models.InterfaceUserService
	logService  models.InterfaceLogService
}

func NewHandler() Handler {

	userData := db.NewUserDBConn("user")
	user := services.NewUserService(userData)

	logData := db.NewTradeLogDBConn("tradelogs")
	log := services.NewLogService(logData)

	return Handler{
		userService: user,
		logService:  log,
	}
}

func (h *Handler) UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return services.ApiResponse(http.StatusMethodNotAllowed, api_error.NewServiceUnavailable())
}

func (h *Handler) HealthCheck(path string) (*events.APIGatewayProxyResponse, error) {
	return services.ApiResponse(http.StatusOK, path)
}

func (h *Handler) GetCookieByName(cookieName string, rawCookie string) (string, error) {
	r := regexp.MustCompile(cookieName + `=\s*(.*?)\s*; `)
	match := r.FindStringSubmatch(rawCookie)
	if len(match) > 0 {
		return match[1], nil
	} else {
		return "", api_error.NewBadRequest("Bad payload")
	}
}
