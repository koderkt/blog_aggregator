package auth

import (
	"errors"
	"net/http"
	"strings"
)



func GetApiKey(header http.Header) (string, error){
	authHeader := header.Get("Authorization")
	if authHeader == ""{
		return "", errors.New("no auth header")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey"{
		return "", errors.New("invalid auth header")
	}

	return splitAuth[1], nil
}
