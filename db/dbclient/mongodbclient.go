package dbclient

import (
	"context"
	"fmt"
	"log"
	"os"

	authaccess "github.com/alubhorta/goth/db/access/auth"
	useraccess "github.com/alubhorta/goth/db/access/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbClient struct {
	_client    *mongo.Client
	UserAccess *useraccess.UserAccess
	AuthAccess *authaccess.AuthAccess
}

func (dbClient *MongoDbClient) Init() {
	log.Println("connecting to db...")

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/admin?w=majority", dbUser, dbPass, dbHost, dbPort)

	ctx := context.Background()
	_mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	dbName := os.Getenv("DB_NAME")
	db := _mongoclient.Database(dbName)

	dbClient._client = _mongoclient
	dbClient.UserAccess = &useraccess.UserAccess{Collection: db.Collection("user")}
	dbClient.AuthAccess = &authaccess.AuthAccess{Collection: db.Collection("userAuthCredential")}

	if err := dbClient._client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln(err)
	}
	log.Println("successfully connected and pinged mongodb! :)")

	// TODO: add necessary indices
}

func (dbClient *MongoDbClient) Cleanup(dbCtx context.Context) {
	log.Println("running DB cleanup...")

	if err := dbClient._client.Disconnect(dbCtx); err != nil {
		log.Fatalln(err)
	}
	dbCtx.Done()
}
