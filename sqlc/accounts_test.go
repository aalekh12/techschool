package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func CreateOwnerAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.GenerateUser(),
		Balance:  util.GenerateMoney(),
		Currency: util.GenerateCurrency(),
	}

	account, err := testQuries.CreateAccount(context.Background(), arg)
	if err != nil {
		fmt.Println("Error in Creating Account", err)
	}

	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateOwnerAccount(t)
}

func TestGetdata(t *testing.T) {
	accounts1 := CreateOwnerAccount(t)
	accounts2, err := testQuries.GetAccount(context.Background(), accounts1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accounts2)

	require.Equal(t, accounts1.ID, accounts2.ID)
	require.Equal(t, accounts1.Owner, accounts2.Owner)
	require.Equal(t, accounts1.Balance, accounts2.Balance)
	require.Equal(t, accounts1.Currency, accounts2.Currency)
	require.WithinDuration(t, accounts1.CreatedAt, accounts2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	accounts1 := CreateOwnerAccount(t)
	arg := UpdateAccountParams{
		ID:      accounts1.ID,
		Balance: util.GenerateMoney(),
	}
	accounts2, err := testQuries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts2)

	require.Equal(t, accounts1.ID, accounts2.ID)
	require.Equal(t, accounts1.Owner, accounts2.Owner)
	require.Equal(t, arg.Balance, accounts2.Balance)
	require.Equal(t, accounts1.Currency, accounts2.Currency)
	require.WithinDuration(t, accounts1.CreatedAt, accounts2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accounts1 := CreateOwnerAccount(t)
	err := testQuries.DeleteAccount(context.Background(), accounts1.ID)
	require.NoError(t, err)

	accounts2, err2 := testQuries.GetAccount(context.Background(), accounts1.ID)
	require.Error(t, err2)
	require.EqualError(t, err2, sql.ErrNoRows.Error())
	require.Empty(t, accounts2)
}

func TestListAccount(t *testing.T) {
	for i := 1; i <= 10; i++ {
		CreateOwnerAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQuries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, v := range accounts {
		require.NotEmpty(t, v)
	}
}
