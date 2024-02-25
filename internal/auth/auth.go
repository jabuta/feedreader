package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var ErrTokenInvalid = errors.New("token invalid")

func GetHeaderAuth(r *http.Request) (string, error) {
	log.Print()
	bearerHeader := strings.Fields(r.Header.Get("Authorization"))
	log.Print(bearerHeader)
	if len(bearerHeader) != 2 || bearerHeader[0] != "ApiKey" {
		return "", ErrTokenInvalid
	}
	return bearerHeader[1], nil
}
