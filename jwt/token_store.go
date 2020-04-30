package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

//token store
type TokenStore interface {
	Save(c jwt.MapClaims) (err error)
	Discard(c jwt.MapClaims) (err error)
}
