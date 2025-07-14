package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(signingMethod jwt.SigningMethod, jwtPayload jwt.MapClaims, signingKey string) (string, error) {
	token := jwt.NewWithClaims(signingMethod, jwtPayload)
	jwtToken, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
