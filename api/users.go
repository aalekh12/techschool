package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/util"
)

type CreateUserAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChnagedAt time.Time `json:"password_chnaged_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) CreateUserAccount(ctx *gin.Context) {
	var req CreateUserAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse(err))
		return
	}
	hashpassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashpassword,
		Email:          req.Email,
		FullName:       req.FullName,
	}

	user, err := server.Store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse(err))
		return
	}

	res := CreateUserResponse{
		FullName:          user.FullName,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChnagedAt: user.PasswordChnagedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, res)
}
