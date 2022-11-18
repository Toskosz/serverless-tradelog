package main

import (
	"github.com/Toskosz/serverless-tradelog/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error) {

	h := handlers.NewHandler()

	if req.HTTPMethod == "GET" && req.Path == "/healthcheck" {
		return h.HealthCheck(req.Path)
	}

	// AUTH ENDPOINTS
	if req.HTTPMethod == "POST" && req.Path == "/register" {
		return h.Register(req)
	} else if req.HTTPMethod == "POST" && req.Path == "/login" {
		return h.Login(req)
	} else if req.HTTPMethod == "POST" && req.Path == "/logout" {
		return h.Logout(req)

		// LOG ENDPOINTS
	} else if req.HTTPMethod == "GET" && req.Path == "/my/logs" {
		if _, ok := req.QueryStringParameters["log-abertura"]; ok {
			return h.GetLog(req)
		} else {
			return h.GetMyLogs(req)
		}
	} else if req.HTTPMethod == "POST" && req.Path == "/my/logs" {
		return h.CreateLog(req)
	} else if req.HTTPMethod == "PUT" && req.Path == "/my/logs" {
		return h.UpdateLog(req)
	} else if req.HTTPMethod == "DELETE" && req.Path == "/my/logs" {
		return h.DeleteLog(req)
	} else {
		return h.UnhandledMethod()
	}
}

func main() {
	lambda.Start(handler)
}

// else if req.HTTPMethod == "GET" && req.Path == "/my/logs" {
//	return h.GetMyLogs(req)
//}
