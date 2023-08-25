package token

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func TestPasteoMaker(t *testing.T) {
	ntoken := strings.TrimSpace(util.RandomString(31))
	log.Println("notk", ntoken)
	maker, err := NewPasteoMaker(ntoken)
	require.NoError(t, err)

	username := util.RandomString(5)
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

func TestExpiredPasteoToken(t *testing.T) {
	make, err := NewPasteoMaker(util.RandomString(31))
	require.NoError(t, err)

	token, err := make.CreateToken(util.GenerateUser(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := make.VerifyToken(token)
	require.Error(t, err)

	require.EqualError(t, err, ErrorExpired.Error())
	require.Nil(t, payload)

}
