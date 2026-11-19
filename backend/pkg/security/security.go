package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("token tidak valid atau sudah kedaluwarsa")
)

type JWTClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// HashPIN mengenkripsi PIN teks biasa menggunakan algoritma bcrypt
func HashPIN(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPINHash memverifikasi kecocokan PIN teks biasa dengan hash tersimpan
func CheckPINHash(pin, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))
	return err == nil
}

// GenerateToken membuat JWT Token berisi claim role pengguna
func GenerateToken(role string, secretKey string, duration time.Duration) (string, int64, error) {
	expTime := time.Now().Add(duration)
	claims := &JWTClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expTime.Unix(), nil
}

// ValidateToken memvalidasi dan mem-parse JWT Token
func ValidateToken(tokenString string, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}