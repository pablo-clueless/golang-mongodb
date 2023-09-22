package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	JwtSecretKey     = []byte(os.Getenv("JWT_SECRET"))
	JwtSigningMethod = jwt.SigningMethodHS256.Name
)

func EncryptPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func VerifyPassword(hashed, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func NewToken(Id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        Id,
		Issuer:    Id,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 860000).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JwtSecretKey, nil
}

func VerifyToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid auth token")
	}

	return claims, nil
}
