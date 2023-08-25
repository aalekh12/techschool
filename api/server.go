package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/token"
	"github.com/techschool/samplebank/util"
)

type Server struct {
	Store      db.Store
	router     *gin.Engine
	tokenmaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenmaker, err := token.NewPasteoMaker(config.TokenSymmetrickey)
	if err != nil {
		return nil, err
	}

	server := &Server{Store: store, config: config, tokenmaker: tokenmaker}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validcurrency)
	}

	server.SrtupServer()
	return server, nil
}

func (server *Server) SrtupServer() {
	router := gin.Default()

	router.POST("/user", server.CreateUserAccount)
	router.POST("/user/login", server.loginUser)
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.getaccount)
	router.GET("/accounts", server.listaccount)
	router.POST("/transfer", server.CreateTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorresponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
