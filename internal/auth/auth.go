package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrTokenInvalid = errors.New("token invalid")

func GetHeaderAuth(r *http.Request) (string, error) {
	bearerHeader := strings.Fields(r.Header.Get("Authorization"))
	if len(bearerHeader) != 2 || bearerHeader[0] != "ApiKey" {
		return "", ErrTokenInvalid
	}
	return bearerHeader[1], nil
}
