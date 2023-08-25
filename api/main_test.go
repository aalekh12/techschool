package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTimeDuration: time.Minute,
		TokenSymmetrickey:  util.RandomString(31),
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
