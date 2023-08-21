package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func CreateRandomTransfers(t *testing.T, account1 Account, account2 Account) Transfer {
	args := CreateTransferesParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Ammount:       util.GenerateMoney(),
	}

	transfer, err := testQuries.CreateTransferes(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Ammount, args.Ammount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCretaeAccout(t *testing.T) {
	ac1 := CreateOwnerAccount(t)
	ac2 := CreateOwnerAccount(t)

	CreateRandomTransfers(t, ac1, ac2)
}

func TestGetTransfers(t *testing.T) {
	ac1 := CreateOwnerAccount(t)
	ac2 := CreateOwnerAccount(t)

	tranfer := CreateRandomTransfers(t, ac1, ac2)

	maketransfer, err := testQuries.GetTransfer(context.Background(), tranfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, maketransfer)

	require.Equal(t, tranfer.Ammount, maketransfer.Ammount)
	require.Equal(t, tranfer.FromAccountID, maketransfer.FromAccountID)
	require.Equal(t, tranfer.ToAccountID, maketransfer.ToAccountID)
	require.Equal(t, tranfer.ID, maketransfer.ID)

}

func TestLsitTransfers(t *testing.T) {
	ac1 := CreateOwnerAccount(t)
	ac2 := CreateOwnerAccount(t)

	for i := 0; i < 5; i++ {

		CreateRandomTransfers(t, ac1, ac2)

		CreateRandomTransfers(t, ac1, ac2)
	}

	args := ListTransferesParams{
		FromAccountID: ac1.ID,
		ToAccountID:   ac2.ID,
		Limit:         5,
		Offset:        5,
	}

	listTransfers, err := testQuries.ListTransferes(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, listTransfers, 5)

	for _, value := range listTransfers {
		require.NotEmpty(t, listTransfers)
		require.True(t, value.FromAccountID == ac1.ID || value.ToAccountID == ac2.ID)
	}

}
