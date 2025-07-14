package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJwt(jwtString string) (jwt.Claims, error) {
	//get the secret key
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		return nil, errors.New("jwt secret key is empty in env")
	}

	//now parse the jwt
	jwtToken, err := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, //this jwt.withValidateMethod will validate method for better protection
		jwt.WithValidMethods([]string{"SHA256"}))

	if err != nil {
		return nil, err
	}

	//now validate against expiry
	// Extract claims if the token is valid
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("failed to extract claims from jwt")
	}

}
