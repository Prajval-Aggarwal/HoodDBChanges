package token

import (
	"fmt"
	"main/server/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Generate JWT Token
func GenerateToken(claims model.Claims) (*string, error) {
	//create user claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))

	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

// Decode Token function
func DecodeToken(tokenString string) (*model.Claims, error) {
	claims := &model.Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}

		return []byte(os.Getenv("JWTKEY")), nil
	})
	//	fmt.Println("claims is", claims)

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func CheckExpiration(tokenString string) (*string, error) {

	claims := &model.Claims{}
	parsedToken, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})
	fmt.Println("pasersed token is", parsedToken)
	accessTokenexpirationTime := time.Now().Add(2 * time.Minute)
	if !claims.VerifyExpiresAt(time.Now(), true) {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(accessTokenexpirationTime)
		newTokenString, err := GenerateToken(*claims)
		if err != nil {
			return nil, err
		}

		return newTokenString, nil

	}
	return &tokenString, nil
}
