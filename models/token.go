package models

import "github.com/dgrijalva/jwt-go"

// Token is the JWT with claims
type Token struct {
	jwt.StandardClaims
	UserID int64
}
