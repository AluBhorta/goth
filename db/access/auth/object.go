package authaccess

import (
	"context"
	"log"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	authmodels "github.com/alubhorta/goth/models/auth"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthAccess struct {
	Collection *mongo.Collection
}

func (ac *AuthAccess) CreateNewUserAuthCredential(credential *authmodels.UserAuthCredential) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ac.Collection.InsertOne(ctx, &credential)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("failed insert of user auth credential.", err)
			return customerrors.ErrDuplicateKey
		}
		return err
	}

	log.Printf("created authCred with mongo_id=%v\n ; userId=%v\n", res.InsertedID, credential.UserId)

	return nil
}
