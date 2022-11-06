package services

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body interface{}) (
	*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func ApiResponseWithCookies(status int, cookies []string, body interface{}) (
	*events.APIGatewayProxyResponse, error) {

	var resp events.APIGatewayProxyResponse

	if len(cookies) > 1 {
		resp = events.APIGatewayProxyResponse{
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: map[string][]string{"Set-Cookie": cookies},
		}
	} else if len(cookies) == 1 {
		resp = events.APIGatewayProxyResponse{
			Headers: map[string]string{"Content-Type": "application/json",
				"Set-Cookie": cookies[0]},
		}
	} else {
		return ApiResponse(status, body)
	}

	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}
