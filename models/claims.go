package models

import "github.com/dgrijalva/jwt-go"

// Claims is the JWT claims
type Claims struct {
	jwt.StandardClaims
	UserID int64
	Role   Role
}
