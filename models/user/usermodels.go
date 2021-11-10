package usermodels

import "time"

type UserInfo struct {
	UserId        string    `json:"userId" bson:"_id"`
	Email         string    `json:"email" bson:"email"`
	FirstName     string    `json:"firstName" bson:"firstName"`
	LastName      string    `json:"lastName" bson:"lastName"`
	Bio           string    `json:"bio" bson:"bio"`
	ProfileImgUrl string    `json:"profileImgUrl" bson:"profileImgUrl"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	ModifiedAt    time.Time `json:"modifiedAt" bson:"modifiedAt"`
}

type CreateUserInfoInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserInfoInput struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Bio           string `json:"bio"`
	ProfileImgUrl string `json:"profileImgUrl"`
}
