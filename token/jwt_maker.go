package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minsecretkey = 32

type JWTMaker struct {
	secretkey string
}

func NEWJwtMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minsecretkey {
		return nil, fmt.Errorf("invalid key size:must be at least %d chracters", minsecretkey)
	}
	return &JWTMaker{secretkey}, nil
}

// CreateToken generate a token for verification
func (maker *JWTMaker) CreateToken(username string, durtaion time.Duration) (string, error) {
	payload, err := NewPayload(username, durtaion)
	if err != nil {
		log.Println(err)
		return "", err
	}
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwttoken.SignedString([]byte(maker.secretkey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalid
		}
		return []byte(maker.secretkey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)
	if err != nil {
		log.Println(err)
		veer, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(veer.Inner, ErrorExpired) {
			return nil, ErrorExpired
		}
		return nil, ErrorInvalid
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalid
	}
	return payload, nil
}
