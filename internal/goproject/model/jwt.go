package model

import (
	"github.com/dgrijalva/jwt-go"
)

type UserJWTClaims struct {
	jwt.StandardClaims

	UID uint32
}
