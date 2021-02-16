package rtoken

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"

	"github.com/birukbelay/item/entity"
)

// CustomClaims specifies custom claims
type CustomClaims struct {
	SessionUUID string
	User        entity.User
	UserId      primitive.ObjectID `json:"_id"`
	jwt.StandardClaims
}

// Categoryte generates jwt token
func Categoryte(signingKey []byte, claims jwt.Claims) (string, error) {
	tn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := tn.SignedString(signingKey)
	return signedString, err
}


// Valid validates a given token
func Valid(signedToken string, signingKey []byte) (*CustomClaims, bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, &CustomClaims{},
	func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil,false, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, true, nil
	}

	return nil, false, err
}

// Claims returns claims used for generating jwt tokens
func Claims(user *entity.User, SessionId string, id primitive.ObjectID, tokenExpires int64) jwt.Claims {
	return CustomClaims{
		SessionId,
		*user ,
		id,
		jwt.StandardClaims{
			ExpiresAt: tokenExpires,
		},
	}
}

// CSRFToken Categorytes random string for CSRF
func CSRFToken(signingKey []byte) (string, error) {
	tn := jwt.New(jwt.SigningMethodHS256)
	signedString, err := tn.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// ValidCSRF checks if a given csrf token is valid
func ValidCSRF(signedToken string, signingKey []byte) (bool, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}
