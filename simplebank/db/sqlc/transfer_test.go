package db

import (
	"context"
	"testing"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

type CreateRandomTransferParams struct {
	FromAccount *Account
	ToAccount   *Account
}

func createRandomTransfer(t *testing.T, params CreateRandomTransferParams) Transfer {
	var account1, account2 Account

	if params.FromAccount != nil {
		account1 = *params.FromAccount
	} else {
		account1 = createRandomAccount(t)
	}

	if params.ToAccount != nil {
		account2 = *params.ToAccount
	} else {
		account2 = createRandomAccount(t)
	}

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, CreateRandomTransferParams{})
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t, CreateRandomTransferParams{})

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)

	require.NotZero(t, transfer2.ID)
	require.NotZero(t, transfer2.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	params := CreateRandomTransferParams{
		FromAccount: &account1,
		ToAccount:   &account2,
	}

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, params)
	}

	arg := ListTransfersParams{
		FromAccountID: params.FromAccount.ID,
		ToAccountID:   params.ToAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
