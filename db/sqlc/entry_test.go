package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/samplebank/util"
)

func CreateEnrty(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Ammount:   util.GenerateMoney(),
	}
	entry, err := testQuries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, entry.Ammount, args.Ammount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

func TestCreateRandomEntry(t *testing.T) {
	myac := CreateOwnerAccount(t)
	CreateEnrty(t, myac)
}

func TestGetEntry(t *testing.T) {
	account := CreateOwnerAccount(t)
	enrty1 := CreateEnrty(t, account)
	entry2, err := testQuries.GetEntry(context.Background(), enrty1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, enrty1.AccountID, entry2.AccountID)
	require.Equal(t, enrty1.Ammount, entry2.Ammount)

	require.WithinDuration(t, entry2.CreatedAt, enrty1.CreatedAt, time.Second)
}

func TestListenrty(t *testing.T) {
	account := CreateOwnerAccount(t)
	for i := 0; i <= 10; i++ {
		CreateEnrty(t, account)
	}

	arge := ListEntryParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQuries.ListEntry(context.Background(), arge)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, allentry := range entries {
		require.NotEmpty(t, allentry)
		require.Equal(t, arge.AccountID, allentry.AccountID)
	}
}
