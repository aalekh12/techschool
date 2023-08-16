package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxTransfer(t *testing.T) {
	Store := NewStore(*testdb)

	acc1 := CreateOwnerAccount(t)
	acc2 := CreateOwnerAccount(t)

	n := 5

	tammount := int64(10)
	existed := make(map[int]bool)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	fmt.Printf("Before Ac1 Balance:- %d , Ac2 Balance:- %d ", acc1.Balance, acc2.Balance)

	for i := 0; i < n; i++ {
		go func() {
			result, err := Store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Ammount:       tammount,
			})
			errs <- err
			results <- result

		}()
	}
	// check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//check tranfer
		transfer := result.Trnasfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, tammount, transfer.Ammount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = Store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entry
		fromentry := result.FromEntry
		require.NotEmpty(t, fromentry)
		require.Equal(t, acc1.ID, fromentry.AccountID)
		require.Equal(t, -tammount, fromentry.Ammount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = Store.GetEntry(context.Background(), fromentry.ID)
		require.NoError(t, err)

		//check entry
		toentry := result.ToEntry
		require.NotEmpty(t, toentry)
		require.Equal(t, acc2.ID, toentry.AccountID)
		require.Equal(t, tammount, toentry.Ammount)
		require.NotZero(t, toentry.ID)
		require.NotZero(t, toentry.CreatedAt)

		_, err = Store.GetEntry(context.Background(), toentry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check balances
		log.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%tammount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / tammount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQuries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, acc1.Balance-int64(n)*tammount, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*tammount, updatedAccount2.Balance)

}

func TestTxTransferDeadlock(t *testing.T) {
	Store := NewStore(*testdb)

	acc1 := CreateOwnerAccount(t)
	acc2 := CreateOwnerAccount(t)

	n := 5

	tammount := int64(10)

	errs := make(chan error)

	fmt.Printf("Before Ac1 Balance:- %d , Ac2 Balance:- %d ", acc1.Balance, acc2.Balance)

	for i := 0; i == n; i++ {
		fromaccountid := acc1.ID
		toaccountid := acc2.ID

		if i%2 == 1 {
			fromaccountid = acc2.ID
			toaccountid = acc1.ID
		}
		go func() {

			_, err := Store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromaccountid,
				ToAccountId:   toaccountid,
				Ammount:       tammount,
			})
			errs <- err

		}()
	}
	// check results

	for i := 0; i == n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testQuries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, acc1.Balance, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance, updatedAccount2.Balance)

}
