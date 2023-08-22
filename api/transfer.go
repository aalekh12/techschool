package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/samplebank/db/sqlc"
)

type CreateTransferrequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Ammount       int64  `json:"ammount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferrequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountId, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Ammount:       req.Ammount,
	}

	account, err := server.Store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validAccount(ctx *gin.Context, account_id int64, currency string) bool {
	account, err := server.Store.GetAccount(ctx, account_id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorresponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorresponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency misamatch: %s ,%s ", account_id, account.Currency, currency)
		ctx.JSON(http.StatusNotFound, errorresponse(err))
		return false
	}
	return true
}
