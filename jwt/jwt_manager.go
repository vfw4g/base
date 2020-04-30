package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vfw4g/base/errors"
	"time"
)

type Config struct {
	SigningMethod string
	Expires       time.Duration
	KeyFunc       func(*jwt.Token) (interface{}, error)
	Issuer        string
	Renewal       bool
	TokenStore    TokenStore
}

type jwtManager struct {
	Config
}

//生成一个jwt_token
func NewManager(c Config) *jwtManager {
	return &jwtManager{c}
}

func (jm *jwtManager) newToken(s interface{}) *jwt.Token {
	now := time.Now()
	newClaims := jwt.MapClaims{
		//"Id":        s.GetId(),
		"ExpiresAt": now.Add(jm.Expires).Unix(),
		"Issuer":    jm.Issuer,
		"IssuedAt":  now.Unix(),
		"subject":   s,
	}
	return jwt.NewWithClaims(jwt.GetSigningMethod(jm.SigningMethod), newClaims)
}

func (jm *jwtManager) valid() error {
	return nil
}

func (jm *jwtManager) NewToken(s interface{}) (token string, err error) {
	var key interface{}
	t := jm.newToken(s)
	if key, err = jm.KeyFunc(t); err != nil {
		return "", errors.WrapMark(err, "keyFunc call error")
	} else {
		return t.SignedString(key)
	}
}

func (jm *jwtManager) NewAndStoreToken(s interface{}) (token string, err error) {
	var key interface{}
	t := jm.newToken(s)
	mc := t.Claims.(jwt.MapClaims)
	if err = jm.TokenStore.Save(mc); err != nil {
		return "", errors.WrapMark(err, "store token error")
	}
	if key, err = jm.KeyFunc(t); err != nil {
		return "", errors.WrapMark(err, "keyFunc call error")
	} else {
		if token, err = t.SignedString(key); err != nil {
			return token, errors.WrapMark(err, "SignedString error")
		} else {
			return
		}
	}
}

//token 失效
func (jm *jwtManager) DiscardToken(c map[string]interface{}) {
	jm.TokenStore.Discard(c)
}
