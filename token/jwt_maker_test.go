package token

import (
	"log"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NEWJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.GenerateUser()
	duration := time.Minute

	issuedat := time.Now()
	expiredat := issuedat.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotZero(t, token)
	log.Println("token =>", token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotZero(t, payload)

	require.NotZero(t, payload.Id)
	require.WithinDuration(t, issuedat, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredat, payload.ExpiredAt, time.Second)

}

func TestExpiredJwtToken(t *testing.T) {
	make, err := NEWJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := make.CreateToken(util.GenerateUser(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := make.VerifyToken(token)
	require.Error(t, err)

	require.EqualError(t, err, ErrorExpired.Error())
	require.Nil(t, payload)

}

func TestInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.GenerateUser(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	make, err := NEWJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = make.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalid.Error())
	require.Nil(t, payload)
}
