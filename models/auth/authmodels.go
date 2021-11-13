package authmodels

import "time"

type UserAuthCredential struct {
	UserId         string    `json:"userId" bson:"_id"`
	Email          string    `json:"email" bson:"email"`
	HashedPassword string    `json:"hashedPassword" bson:"hashedPassword"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	ModifiedAt     time.Time `json:"modifiedAt" bson:"modifiedAt"`
}

type SignupInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutInput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshInput struct {
	RefreshToken string `json:"refreshToken"`
}

type ResetInitInput struct {
	Email string `json:"email"`
}

type ResetVerifyInput struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"newPassword"`
}
