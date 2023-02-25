package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint32
	NickName    string
	AuthorityId uint32
	jwt.StandardClaims
}
