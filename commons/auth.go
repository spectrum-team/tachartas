package commons

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spectrum-team/tachartas/models"
)

var secret = []byte("el-famoso-biberon")

const ctxKey = "secretctxkey"

func SignIn(user *models.User) string {

	expires := time.Now().Add(time.Minute * 60)
	c := models.AuthInfo{
		Email:    user.Email,
		FullName: user.DisplayName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenStr, _ := token.SignedString(secret)

	return tokenStr
}

func AuthMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		ctx := r.Context()

		if token == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := &models.AuthInfo{}

		t, err := jwt.ParseWithClaims(token, c, func(tkn *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.Header().Set("WWW-Authenticate", "Bearer")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !t.Valid {
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, ctxKey, c.Email)
		w.Header().Set("Authorization", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(f)
}

func GetAuthCtx(ctx context.Context) string {
	if v := ctx.Value(ctxKey); v != nil {
		return v.(string)
	}

	return ""
}
