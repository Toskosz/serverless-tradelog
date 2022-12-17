package services

import (
	"net/http"

	"github.com/Toskosz/serverless-tradelog/models/api_error"
)

func GetCookieByName(cookieName string, rawCookies string) (string, error) {

	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return "", api_error.NewBadRequest(rawCookies)
	}

	return cookie.Value, nil
}
