package db

import (
	"context"
	"testing"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

type CreateRandomEntryParams struct {
	Account *Account
}

func createRandomEntry(t *testing.T, params CreateRandomEntryParams) Entry {
	var account Account

	if params.Account != nil {
		account = *params.Account
	} else {
		account = createRandomAccount(t)
	}

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, CreateRandomEntryParams{})
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t, CreateRandomEntryParams{})

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.Amount, entry1.Amount)
	require.NotZero(t, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, CreateRandomEntryParams{Account: &account})
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
