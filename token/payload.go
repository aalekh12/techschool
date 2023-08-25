package token

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorExpired = errors.New("Token is Expired")
	ErrorInvalid = errors.New("Token is Invalid")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenid, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	payload := &Payload{
		Id:        tokenid,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil

}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpired
	}
	return nil
}
