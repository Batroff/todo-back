package middleware

import (
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"
)

func Auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if ok, err := IsTokenValid(r); err != nil || !ok {
		log.Printf("ok %v, err %v", ok, err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(rw, r)
}

func CreateToken(id models.ID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func ParseToken(r *http.Request) (*jwt.Token, error) {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("secret")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func IsTokenValid(r *http.Request) (bool, error) {
	token, err := ParseToken(r)
	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return false, err
	}
	return true, nil
}
