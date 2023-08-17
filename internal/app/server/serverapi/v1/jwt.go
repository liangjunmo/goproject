package v1

import (
	"github.com/dgrijalva/jwt-go"
)

type UserJwtClaims struct {
	jwt.StandardClaims

	Uid uint32
}
