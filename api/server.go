package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/samplebank/db/sqlc"
)

type Server struct {
	Store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{Store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validcurrency)
	}

	router.POST("/user", server.CreateUserAccount)
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.getaccount)
	router.GET("/accounts", server.listaccount)
	router.POST("/transfer", server.CreateTransfer)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorresponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
