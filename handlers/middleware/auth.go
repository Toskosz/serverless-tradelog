package middleware

import (
	"net/http"
	"os"

	"github.com/Toskosz/serverless-tradelog/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type authResponse struct {
	message string
}

func Auth(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	tokenString, err := services.GetCookieByName("jwt", req.Headers["Cookie"])
	if err != nil {
		return services.ApiResponse(http.StatusUnauthorized, authResponse{"Unauthenticated"})
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// check token signing method etc
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return services.ApiResponse(http.StatusUnauthorized, authResponse{"Unauthenticated"})
	}

	if _, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return services.ApiResponse(http.StatusOK, authResponse{"sucess"})
	} else {
		return services.ApiResponse(http.StatusUnauthorized, authResponse{"Unauthenticated"})
	}
}
