package db

import (
	"context"
	"fmt"
	"time"
)

// TransferTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

// TransferTxResult is the result of the transfer transaction
type CreateUserTxResult struct {
	UserID int64 `json:"user_id"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		startTime := time.Now()
		user, err := q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			HashedPassword: arg.PasswordHash,
			Email:          arg.Email,
		})

		if err != nil {
			return err
		}
		fmt.Println("Created user time =>", time.Now().Sub(startTime))

		result.UserID = user.ID

		return err

	})

	return result, err
}
