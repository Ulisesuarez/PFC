package middleware

import (
	"bitbucket.org/inehealth/api/api"
	"bitbucket.org/inehealth/idonia-pacs/service/auth"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type Authorization struct {
}

func (m Authorization) Name() string {
	return "authorization"
}

func (m Authorization) Middleware() api.RESTMiddlewareFunction {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authorization := &auth.Auth{}
		authorizationHeader := r.Header.Get("Authorization")
		switch {
		case strings.HasPrefix(authorizationHeader, "IDONIA "):
			authorization.Token = strings.TrimPrefix(authorizationHeader, "IDONIA ")
		case strings.HasPrefix(authorizationHeader, "Bearer"):
			bearer := strings.Trim(authorizationHeader[len("Bearer"):], " ")
			token, _ := jwt.ParseWithClaims(bearer, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return "", nil //TODO: Check signature
			})
			authorization.Token = token.Claims.(*jwt.StandardClaims).Subject
		default:
		}
		ctx := context.WithValue(r.Context(), "authorization", authorization)

		next(rw, r.WithContext(ctx))
	}
}
