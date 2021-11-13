package useraccess

import (
	"context"
	"log"
	"time"

	customerrors "github.com/alubhorta/goth/custom/errors"
	usermodels "github.com/alubhorta/goth/models/user"

	"go.mongodb.org/mongo-driver/bson"
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

func (ac *UserAccess) GetAUser(userId string) (*usermodels.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo := new(usermodels.UserInfo)
	result := ac.Collection.FindOne(ctx, bson.M{"_id": userId})
	err := result.Decode(userInfo)
	if err == mongo.ErrNoDocuments {
		return nil, customerrors.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (ac *UserAccess) UpdateAUser(userId string, input *usermodels.UpdateUserInfoInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := ac.Collection.UpdateOne(
		ctx,
		bson.M{"_id": userId},
		bson.D{
			{"$set", bson.D{
				{Key: "firstName", Value: input.FirstName},
				{Key: "lastName", Value: input.LastName},
				{Key: "bio", Value: input.Bio},
				{Key: "profileImgUrl", Value: input.ProfileImgUrl},
				{Key: "modifiedAt", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("failed update of user.", err)
			return customerrors.ErrDuplicateKey
		}
		return err
	} else if result.MatchedCount == 0 {
		return customerrors.ErrNotFound
	}

	return nil
}

func (ac *UserAccess) DeleteAUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := ac.Collection.DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		return err
	} else if result.DeletedCount == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}
