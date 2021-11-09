package commonclients

import "github.com/alubhorta/goth/db/dbclient"

type CommonClients struct {
	DbClient *dbclient.MongoDbClient
}
