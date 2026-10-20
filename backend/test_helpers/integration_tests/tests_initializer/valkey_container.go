package tests_initializer

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/valkey"
)

type VkContainer struct {
	container *valkey.ValkeyContainer
	connStr   string
	rdb       *redis.Client
}

func SetupValkey(ctx context.Context) (*VkContainer, error) {
	valkeyContainer, err := valkey.Run(ctx,
		"valkey/valkey:7.2.5",
		valkey.WithSnapshotting(10, 1),
		valkey.WithLogLevel(valkey.LogLevelVerbose),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	connStr, err := valkeyContainer.ConnectionString(ctx)
	if err != nil {
		log.Printf("failed to get connection string: %s", err)
		return nil, err
	}

	opts, err := redis.ParseURL(connStr)
	if err != nil {
		log.Printf("failed to parse connection string: %s", err)
		return nil, err
	}

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("failed to ping valkey: %s", err)
		return nil, err
	}

	return &VkContainer{
		container: valkeyContainer,
		connStr:   connStr,
		rdb:       rdb,
	}, nil
}

func (vc *VkContainer) Teardown(ctx context.Context) error {
	if vc.rdb != nil {
		if err := vc.rdb.Close(); err != nil {
			log.Printf("failed to close redis client: %s", err)
			return err
		}
	}
	if vc.container != nil {
		if err := vc.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate valkey container: %s", err)
			return err
		}
	}
	return nil
}
