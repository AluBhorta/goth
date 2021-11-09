package useraccess

import (
	"context"
	"log"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	usermodels "github.com/alubhorta/goth/models/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserAccess struct {
	Collection *mongo.Collection
}

func (ac *UserAccess) CreateAUser(userId string, input *usermodels.CreateUserInfoInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	userInfo := usermodels.UserInfo{
		UserId:        userId,
		Email:         input.Email,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Bio:           "",
		ProfileImgUrl: "",
		CreatedAt:     now,
		ModifiedAt:    now,
	}
	res, err := ac.Collection.InsertOne(ctx, &userInfo)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("failed insert of user.", err)
			return customerrors.ErrDuplicateKey
		}
		return err
	}

	log.Printf("created user with mongo_id=%v\n ; userId=%v\n", res.InsertedID, userId)

	return nil
}
