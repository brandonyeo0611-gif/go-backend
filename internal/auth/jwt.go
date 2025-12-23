package auth

import (
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
// create a claim structure to retreive the userid 

func GenerateToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID: userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims) 
	// creates a token with the claims
	// signs the token using 
	return token.SignedString(([]byte("tryingtoauth")))
	}

	func ValidateToken(tokenStr string) (*Claims, error) {
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token*jwt.Token) (interface{}, error){
			return []byte("tryingtoauth"),nil
		})
		// verify signature and fills the empty claim struc with decoded data and return token

		if err != nil {
			return nil, err
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid{
			return claims, nil
		}
		// return userid and uername
		return nil, fmt.Errorf("invalid token")
	}