package api

import (
	"database/sql"
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

type loginuserAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginresponse struct {
	AccessToken string             `json:"access_token"`
	User        CreateUserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginuserAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse(err))
		return
	}

	user, err := server.Store.GetUsers(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorresponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorresponse(err))
		return
	}
	err = util.ComparePassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorresponse(err))
		return
	}

	accesstoken, err := server.tokenmaker.CreateToken(user.Username, server.config.AccessTimeDuration)
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
	rsp := loginresponse{
		AccessToken: accesstoken,
		User:        res,
	}

	ctx.JSON(http.StatusOK, rsp)
}
