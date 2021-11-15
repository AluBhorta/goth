package transactions

import (
	"context"
	"fmt"

	customerrors "github.com/alubhorta/goth/custom/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DbTxnManager struct {
	Client      *mongo.Client
	AuthCredCol *mongo.Collection
	UserCol     *mongo.Collection
}

func (dtm *DbTxnManager) DeleteUserTxn(userId string) error {
	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.
		if result, err := dtm.AuthCredCol.DeleteOne(sessCtx, bson.M{"_id": userId}); err != nil {
			return nil, err
		} else if result.DeletedCount == 0 {
			return nil, customerrors.ErrNotFound
		}

		if result, err := dtm.UserCol.DeleteOne(sessCtx, bson.M{"_id": userId}); err != nil {
			return nil, err
		} else if result.DeletedCount == 0 {
			return nil, customerrors.ErrNotFound
		}

		return nil, nil
	}

	// Step 2: Start a session and run the callback using WithTransaction.
	session, err := dtm.Client.StartSession()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}
	fmt.Printf("txn result: %v\n", result)

	return nil
}
