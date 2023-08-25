package token

import "time"

// Maker interface will hold the functions for creating and matching tokens
type Maker interface {

	//CreateToken generate a token for verification
	CreateToken(username string, durtaion time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
