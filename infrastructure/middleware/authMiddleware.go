package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.URL.Path)
		if r.URL.Path == "/registration" || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			role := "ANONYMOUS"
			isAuthorized, err := Enforce(role, r.URL.Path, r.Method)
			if !isAuthorized {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized request: " + err.Error()))
				next.ServeHTTP(w, r)
				return
			}

			r.Header.Set("role", role)
			next.ServeHTTP(w, r)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}

		username := claims.(jwt.MapClaims)["username"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		isAuthorized, err := Enforce(role, r.URL.Path, r.Method)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error while authorization: " + err.Error()))
			return
		}

		if !isAuthorized {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized request: " + err.Error()))
			next.ServeHTTP(w, r)
			return
		}

		r.Header.Set("username", username)
		r.Header.Set("role", role)

		next.ServeHTTP(w, r)
	})
}

func Enforce(role string, obj string, act string) (bool, error) {
	m, _ := os.Getwd()
	fmt.Println(m)

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", "config/rbac_policy.csv")
	if err != nil {
		return false, fmt.Errorf("failed to create enforcer: %w", err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy: %w", err)
	}
	ok, _ := enforcer.Enforce(role, obj, act)
	return ok, nil
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("123456")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
