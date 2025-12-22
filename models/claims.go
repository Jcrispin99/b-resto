package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}
