package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/common/logger"
	"github.com/songquanpeng/one-api/controller"
	"github.com/songquanpeng/one-api/middleware"
	"github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/router"
)

func main() {
	// Print version information
	common.PrintVersion()

	// Load configuration from environment
	config.LoadConfig()

	// Initialize logger
	if err := logger.SetupLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	logger.SysLog(fmt.Sprintf("New API %s starting...", common.Version))

	// Initialize database
	if err := model.InitDB(); err != nil {
		logger.FatalLog(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer model.CloseDB()

	// Run database migrations
	if err := model.MigrateDB(); err != nil {
		logger.FatalLog(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	// Initialize Redis if configured
	if config.RedisEnabled {
		if err := common.InitRedisClient(); err != nil {
			// Treat Redis failure as fatal since we rely on it for rate limiting and caching
			logger.FatalLog(fmt.Sprintf("Failed to initialize Redis: %v", err))
		}
	}

	// Initialize default options and admin user
	model.InitOptionMap()
	if err := model.InitRootUser(); err != nil {
		logger.SysError(fmt.Sprintf("Failed to initialize root user: %v", err))
	}

	// Start background tasks
	go model.SyncOptions(config.SyncFrequency)
	go controller.AutomaticDisableChannel()
	go controller.AutomaticEnableChannel()

	// Setup Gin engine
	if config.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	server.Use(middleware.CORS())

	// Register all routes
	router.SetRouter(server)

	// Determine listen address
	addr := fmt.Sprintf(":%d", config.Port)
	logger.SysLog(fmt.Sprintf("Server listening on %s", addr))

	if err := server.Run(addr); err != nil {
		logger.FatalLog(fmt.Sprintf("Failed to start server: %v", err))
	}
}
