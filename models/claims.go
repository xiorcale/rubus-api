package models

import "github.com/dgrijalva/jwt-go"

// Claims is the JWT claims
type Claims struct {
	jwt.StandardClaims
	UserID int64
	Role   Role
}

// JWT is used for swagger doc only
type JWT struct {
	token string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwic3ViIjoxfQ.ThUA0fhJMGLGtBFAALQ8zdczOzlRIJsV3UY3GKpEZH0"`
}
