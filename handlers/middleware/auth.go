package middleware

import (
	"net/http"
	"os"

	"github.com/Toskosz/serverless-tradelog/models/api_error"
	"github.com/Toskosz/serverless-tradelog/services"
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
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if err != nil {
		return services.ApiResponse(http.StatusUnauthorized,
			api_error.NewAuthorization(err.Error()))
	}

	//if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
	//	return services.ApiResponse(http.StatusOK, authResponse{"sucess"})
	//} else {
	//	return nil, api_error.NewAuthorization("failed to get user from token")
	//}

	return services.ApiResponse(http.StatusOK, authResponse{"sucess"})
}
