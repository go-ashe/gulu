package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func NewClient(ctx context.Context, cs map[string]Config) (map[string]redis.UniversalClient, error) {
	clis := make(map[string]redis.UniversalClient, len(cs))
	for key, conf := range cs {
		var cli redis.UniversalClient
		switch conf.ClientType {
		case "cluster":
			cli = redis.NewClusterClient(&redis.ClusterOptions{
				Username:        conf.Username,
				Password:        conf.Password,
				Addrs:           conf.Addrs,
				PoolSize:        conf.PoolSize,
				DialTimeout:     conf.DialTimeout,
				ReadTimeout:     conf.ReadTimeout,
				WriteTimeout:    conf.WriteTimeout,
				MaxRetries:      conf.MaxRetries,
				MinRetryBackoff: conf.MinRetryBackoff,
				MaxRetryBackoff: conf.MaxRetryBackoff,
			})
		default:
			cli = redis.NewClient(&redis.Options{
				Addr:            conf.Addrs[0],
				Username:        conf.Username,
				Password:        conf.Password,
				DB:              conf.DB,
				PoolSize:        conf.PoolSize,
				DialTimeout:     conf.DialTimeout,
				ReadTimeout:     conf.ReadTimeout,
				WriteTimeout:    conf.WriteTimeout,
				MaxRetries:      conf.MaxRetries,
				MinRetryBackoff: conf.MinRetryBackoff,
				MaxRetryBackoff: conf.MaxRetryBackoff,
			})
		}
		_, err := cli.Ping(ctx).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "failed ping redis %v", key)
		}
		clis[key] = cli
	}
	return clis, nil
}
