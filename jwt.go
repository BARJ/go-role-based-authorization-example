package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	jwt.StandardClaims
	Name  string `json:"name,omitempty"`
	Roles []Role `json:"roles,omitempty"`
}

func encodeJwt(claims JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("encodeJwt: %v", err)
	}
	return tokenString, nil
}

func decodeJwt(s string) (*jwt.Token, *JwtClaims, error) {
	token, err := jwt.ParseWithClaims(s, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("decodeJwt: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		return nil, nil, fmt.Errorf("decodeJwt: %v", err)
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, nil, errors.New("decodeJwt: invalid claims")

	}
	return token, claims, nil
}
