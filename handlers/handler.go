package handlers

import (
	"net/http"

	"github.com/Toskosz/everythingreviewed/db"
	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/Toskosz/everythingreviewed/services"
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
