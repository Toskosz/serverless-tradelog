package main

import (
	"github.com/Toskosz/everythingreviewed/db"
	"github.com/Toskosz/everythingreviewed/handlers"
	"github.com/Toskosz/everythingreviewed/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	userData := db.NewUserDBConn("user")
	userService := services.NewUserService(userData)
	h := handlers.NewHandler(userService)

	if req.HTTPMethod == "POST" && req.Path == "/register" {
		return h.Register(req)
	} else if req.HTTPMethod == "POST" && req.Path == "/login" {
		return h.Login(req)
	} else if req.HTTPMethod == "POST" && req.Path == "/logout" {
		return h.Logout(req)
	} else {
		return h.UnhandledMethod()
	}
}

func main() {
	lambda.Start(handler)
}
