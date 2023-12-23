package accountservice

import (
	"github.com/dgrijalva/jwt-go"
)

type UserJwtClaims struct {
	jwt.StandardClaims

	UID uint32
}
