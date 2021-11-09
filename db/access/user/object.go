package useraccess

import "go.mongodb.org/mongo-driver/mongo"

type UserAccess struct {
	Collection *mongo.Collection
}
