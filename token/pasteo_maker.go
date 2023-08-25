package token

import (
	"fmt"
	"log"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasteoMaker struct {
	pasteo      *paseto.V2
	symetrickey []byte
}

func NewPasteoMaker(symetrickey string) (Maker, error) {
	keyBytes := []byte(symetrickey)
	fmt.Print("size==>", len(symetrickey))
	fmt.Println("\nch key size0", chacha20poly1305.KeySize)
	if len(symetrickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Error invalid key size: must be exactly %d bytes", chacha20poly1305.KeySize)
	}
	maker := &PasteoMaker{
		pasteo:      paseto.NewV2(),
		symetrickey: keyBytes,
	}
	return maker, nil
}

// CreateToken generate a token for verification
func (maker *PasteoMaker) CreateToken(username string, durtaion time.Duration) (string, error) {
	payload, err := NewPayload(username, durtaion)
	if err != nil {
		return "", err
	}
	token, err := maker.pasteo.Encrypt(maker.symetrickey, payload, nil)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	log.Println("pasteo ==>", token)
	return token, nil
}

func (maker *PasteoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.pasteo.Decrypt(token, maker.symetrickey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
