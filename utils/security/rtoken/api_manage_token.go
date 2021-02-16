package rtoken

import (
	"net/http"

"fmt"
	"github.com/dgrijalva/jwt-go"
)

// Create creates and sets session cookie
func CreateToken(w http.ResponseWriter, signingKey []byte, claims jwt.Claims ) {

	signedString, err := Categoryte(signingKey, claims)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Token", signedString)
}

//// Valid validates client cookie value
//func Valid(SignedToken string, signingKey []byte) (bool, error) {
//	valid, err := rtoken2.Valid(SignedToken, signingKey)
//	if err != nil || !valid {
//		return false, errors.New("Invalid Session Cookie")
//	}
//	return true, nil
//}