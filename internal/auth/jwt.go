package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		panic("[pkg auth func GenerateJWT] env: SECRET_KEY is missing")
	}

	return []byte(key)
}

func GenerateJWT(userID string) (string, error) {
	secretKey := getSecretKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "mocker-auth",
		"sub": userID,
		"iat": time.Now().Unix(),
	})

	signed, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("[pkg auth func GenerateJWT] failed to sign :: %w", err)
	}

	return signed, nil
}

func VerifyJWT(token string) (*jwt.Token, jwt.MapClaims, error) {	
	secretKey := getSecretKey()

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[pkg auth func VerifyJWT] unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil 
	})

	if err != nil {
		return nil, nil, err 
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return t, claims, nil
	}
	return nil, nil, fmt.Errorf("unexpected token")
}
