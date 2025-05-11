package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/pawverse/pawcare-core/pkg/data/mongodb"
	"github.com/pawverse/pawcare-profiles/internal/conf"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOwnerRepo)

// Data .
type Data struct {
	db *mongo.Database
}

// NewData .
func NewData(ctx context.Context, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	database, dbCleanup, err := mongodb.ConnectDatabase(ctx, mongodb.ConnectionStringConfig(c.Database), c.Database.Name)
	if err != nil {
		return nil, func() {}, err
	}

	cleanup := func() {
		err := dbCleanup(ctx)
		if err != nil {
			panic(err)
		}
	}

	return &Data{
		db: database,
	}, cleanup, nil
}
