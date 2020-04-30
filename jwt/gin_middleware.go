package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/vfw4g/base/errors"
	"github.com/vfw4g/base/global"
)

var (
	claimsExtractKey = "CUR_CLAIMS"
	logger           = global.Logger
)

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get(claimsExtractKey)
	if !exists {
		return make(jwt.MapClaims)
	}
	return claims.(jwt.MapClaims)
}

// ExtractClaimsFromToken help to extract the JWT claims from token
func extractClaimsFromToken(token *jwt.Token) jwt.MapClaims {
	if token == nil {
		return make(map[string]interface{})
	}
	return token.Claims.(jwt.MapClaims)
}

// GetToken help to get the JWT token string
func GetToken(c *gin.Context) string {
	token, exists := c.Get(claimsExtractKey)
	if !exists {
		return ""
	}

	return token.(string)
}

func JwtTokenHandlerFunc(jm *jwtManager) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		// only accessible with a valid token
		// Get token from request
		token, err := request.ParseFromRequest(
			cxt.Request,
			request.AuthorizationHeaderExtractor,
			jm.KeyFunc,
			request.WithClaims(jwt.MapClaims{}),
		)
		if err != nil {
			logger.Errorvn(errors.WrapMark(err, "token parse error"))
			//临时注释
			//cxt.AbortWithStatusJSON(http.StatusUnauthorized, rsp.StatusBadRequest(err.Error()))
			return
		}
		if err := jm.valid(); err != nil {
			logger.Errorvn(errors.WrapMark(err, "token validate error"))
			//临时注释
			//cxt.AbortWithStatusJSON(http.StatusUnauthorized, rsp.StatusBadRequest(err.Error()))
			return
		}
		// Token is valid
		claims := extractClaimsFromToken(token)
		cxt.Set(claimsExtractKey, claims)
		logger.Infof("%+v\n", claims)
		//cxt.Next()
	}
}
