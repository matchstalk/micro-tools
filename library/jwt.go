package library

import (
	"github.com/matchstalk/jwt"
)

func GenerateJwt(secret string, claims *jwt.Claims) (token string, err error) {
	algorithm := jwt.HmacSha256(secret)
	token, err = algorithm.Encode(claims)
	return token, err
}

//验证jwt
func VerifyJwt(secret, token string) (c *jwt.Claims, err error) {
	algorithm := jwt.HmacSha256(secret)
	return algorithm.DecodeAndValidate(token)
}
