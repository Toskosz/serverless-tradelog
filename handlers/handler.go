package handlers

import (
	"net/http"

	"github.com/Toskosz/everythingreviewed/models"
	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/Toskosz/everythingreviewed/services"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	userService models.InterfaceUserService
}

func NewHandler(user models.InterfaceUserService) Handler {
	return Handler{
		userService: user,
	}
}

func (h *Handler) UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return services.ApiResponse(http.StatusMethodNotAllowed, api_error.NewServiceUnavailable())
}
