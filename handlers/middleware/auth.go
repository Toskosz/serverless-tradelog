package middleware

import (
	"net/http"
	"os"

	"github.com/Toskosz/everythingreviewed/models/api_error"
	"github.com/Toskosz/everythingreviewed/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type authResponse struct {
	message string
}

func Auth(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	token, exists := req.Headers["jwt"]
	if exists {
		return services.ApiResponse(http.StatusUnauthorized, authResponse{"Unauthenticated"})
	}

	_, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

	if err != nil {
		return services.ApiResponse(http.StatusUnauthorized,
			api_error.NewAuthorization(err.Error()))
	}

	return services.ApiResponse(http.StatusOK, authResponse{"sucess"})
}
