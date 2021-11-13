package commonmodels

import (
	"github.com/alubhorta/goth/db/cacheclient"
	"github.com/alubhorta/goth/db/dbclient"
)

type CommonCtx struct {
	Clients *CommonClients
	UserId  string
}

type CommonClients struct {
	DbClient    *dbclient.MongoDbClient
	CacheClient *cacheclient.RedisClient
}
