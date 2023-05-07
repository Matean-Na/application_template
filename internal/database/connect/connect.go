package connect

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB
var RedisDB *redis.Client
