package authaccess

import "go.mongodb.org/mongo-driver/mongo"

type AuthAccess struct {
	Collection *mongo.Collection
}
