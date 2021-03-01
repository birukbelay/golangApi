package rtoken

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// Create creates and sets session cookie
func CreateToken( signingKey []byte, claims jwt.Claims )(string, error ){

	signedString, err := Categoryte(signingKey, claims)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	return signedString, nil


}

//// Valid validates client cookie value
//func Valid(SignedToken string, signingKey []byte) (bool, error) {
//	valid, err := rtoken2.Valid(SignedToken, signingKey)
//	if err != nil || !valid {
//		return false, errors.New("Invalid Session Cookie")
//	}
//	return true, nil
//}