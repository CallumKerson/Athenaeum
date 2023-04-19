package main

import (
	"strings"

	"github.com/CallumKerson/Athenaeum/internal/adapters/bolt"
	"github.com/CallumKerson/Athenaeum/internal/adapters/logrus"
	"github.com/CallumKerson/Athenaeum/internal/memcache"
	"github.com/CallumKerson/Athenaeum/pkg/client"
)

func (c *Config) GetBoltDBOps() []bolt.Option {
	return []bolt.Option{bolt.WithDBDefaults(), bolt.WithPathToDBDirectory(c.DB.Root)}
}

func (c *Config) GetClientOpts() []client.Option {
	return []client.Option{client.WithHost(c.Host), client.WithVersion(Version)}
}

func (c *Config) GetMemcacheOpts() []memcache.Option {
	ttl, _ := c.Cache.GetTTL()
	return []memcache.Option{memcache.WithTTL(ttl), memcache.WithCapacity(c.Cache.Length)}
}

func (c *Config) GetLogger() *logrus.Logger {
	log := logrus.NewLogger()
	level := strings.ToLower(c.Log.Level)
	switch level {
	case "debug":
		log.SetLevelDebug()
	case "warn":
		log.SetLevelWarn()
	case "error":
		log.SetLevelError()
	default:
		log.SetLevelInfo()
	}
	return log
}
