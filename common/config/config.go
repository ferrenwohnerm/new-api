package config

import (
	"os"
	"strconv"
	"sync"
)

var (
	once sync.Once

	// Server configuration
	ServerPort       = 3000
	ServerAddress    = ""
	GinMode          = "debug"

	// Database configuration
	SQLitePath       = "new-api.db"
	MySQLDSN         = ""
	PostgresDSN       = ""

	// Redis configuration
	RedisConnString  = ""
	RedisPassword    = ""
	RedisMasterName  = ""

	// Security
	SessionSecret    = "new-api-secret"
	CryptoSecret     = ""

	// Feature flags
	DebugEnabled     = false
	DemoMode         = false
	// Enabled by default for personal use — avoids repeated DB hits on a low-resource VPS
	MemoryCacheEnabled = true

	// Rate limiting
	// Increased from 60 to 120 to better suit personal/self-hosted use
	GlobalApiRateLimitNum      = 120
	GlobalApiRateLimitDuration = int64(3 * 60)

	// Model & channel defaults
	DefaultChannelModels    = map[int][]string{}
	DefaultChannelModelMapping = map[int]map[string]string{}

	// System info
	Version = "v0.0.1"
	StartTime int64
)

// InitConfig loads configuration from environment variables.
// It is safe to call multiple times; only the first call takes effect.
func InitConfig() {
	once.Do(func() {
		loadFromEnv()
	})
}

func loadFromEnv() {
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			ServerPort = p
		}
	}

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		ServerAddress = addr
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		GinMode = mode
	}

	if dsn := os.Getenv("SQL_DSN"); dsn != "" {
		MySQLDSN = dsn
	}

	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		PostgresDSN = dsn
	}

	if path := os.Getenv("SQLITE_PATH"); path != "" {
		SQLitePath = path
	}

	if redis := os.Getenv("REDIS_CONN_STRING"); redis != "" {
		RedisConnString = redis
	}

	if redisPwd := os.Getenv("REDIS_PASSWORD"); redisPwd != "" {
		RedisPassword = redisPwd
	}

	if secret := os.Getenv("SESSION_SECRET"); secret != "" {
		SessionSecret = secret
	}

	if crypto := os.Getenv("CRYPTO_SECRET"); crypto != "" {
		CryptoSecret = crypto
	}

	if debug := os.Getenv("DEBUG"); debug == "true" || debug == "1" {
		DebugEnabled = true
	}

	if demo := os.Getenv("DEMO_MODE"); demo == "true" || demo == "1" {
		DemoMode = true
	}

	// Allow explicitly disabling memory cache via env even though it defaults to true
	if cache := os.Getenv("MEMORY_CACHE_ENABLED"); cache == "false" || cache == "0" {
		MemoryCacheEnabled = false
	} else if cache == "true" || cache == "1" {
		MemoryCacheEnabled = true
	}

	if rateLimit := os.Getenv("GLOBAL_API_RATE_LIMIT"); rateLimit != "" {
		if n, err := strconv.Atoi(rateLimit); err == nil {
			GlobalApiRateLimitNum = n
		}
	}
}
