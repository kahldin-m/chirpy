// auth.go
package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// run a password through argon2 and return a hashed version
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		fmt.Printf("Fatal error: %v", err)
	}

	return hash, err
}

// check a password against the stored hashed version
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Printf("Fatal error: %v", err)
	}

	return match, err
}

// make and return a signed JSON Web Token
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now().UTC()
	expiry := now.Add(expiresIn)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		Subject: uuid.UUID.String(userID),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiry),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}



func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		 return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	finalID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return finalID, nil
}