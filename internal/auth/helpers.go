package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedString, err := jwt.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if token.Valid != true {
		return uuid.Nil, errors.New("Token is invalid")
	}
	userId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userUUID, err := uuid.Parse(userId)
	if token.Valid != true {
		return uuid.Nil, errors.New("Invalid sub invalid")
	}
	return userUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeaderParts, s, err := getAuthorizationHeadersParts(headers)
	if err != nil {
		return s, err
	}
	if authorizationHeaderParts[0] != "Bearer" {
		return "", errors.New("Authorization header is invalid")
	}
	return authorizationHeaderParts[1], nil
}

func getAuthorizationHeadersParts(headers http.Header) ([]string, string, error) {
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == "" {
		return nil, "", errors.New("Authorization header is missing")
	}
	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 {
		return nil, "", errors.New("Authorization header is invalid")
	}
	return authorizationHeaderParts, "", nil
}

func GetApiKey(headers http.Header) (string, error) {
	authorizationHeaderParts, s, err := getAuthorizationHeadersParts(headers)
	if err != nil {
		return s, err
	}
	if authorizationHeaderParts[0] != "ApiKey" {
		return "", errors.New("Authorization header is invalid")
	}
	return authorizationHeaderParts[1], nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
