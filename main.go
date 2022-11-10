package main

import (
	"github.com/Toskosz/everythingreviewed/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	h := handlers.NewHandler()

	if req.HTTPMethod == "POST" && req.Resource == "/register" {
		return h.Register(req)
	} else if req.HTTPMethod == "POST" && req.Resource == "/login" {
		return h.Login(req)
	} else if req.HTTPMethod == "POST" && req.Resource == "/logout" {
		return h.Logout(req)
	} else if req.HTTPMethod == "GET" && req.Resource == "/my/logs/{log-ts}" {
		return h.GetLog(req)
	} else if req.HTTPMethod == "GET" && req.Resource == "/my/logs" {
		return h.GetMyLogs(req)
	} else if req.HTTPMethod == "POST" && req.Resource == "/my/logs" {
		return h.CreateLog(req)
	} else if req.HTTPMethod == "PUT" && req.Resource == "/my/logs" {
		return h.UpdateLog(req)
	} else if req.HTTPMethod == "DELETE" && req.Resource == "/my/logs/{log-ts}" {
		return h.DeleteLog(req)
	} else {
		return h.UnhandledMethod()
	}
}

func main() {
	lambda.Start(handler)
}
