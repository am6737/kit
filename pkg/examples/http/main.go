package main

import (
	"context"
	"github.com/cossim/coss-server/pkg/encryption"
	ctrl "github.com/cossim/kit"
	"github.com/cossim/kit/pkg/config"
	"github.com/cossim/kit/pkg/healthz"
	"github.com/cossim/kit/pkg/log"
	"github.com/cossim/kit/pkg/version"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	mgr, err := ctrl.NewManager(nil, ctrl.Options{
		Http: ctrl.HTTPServer{
			//HTTPService: &Test{},
		},
		Grpc: ctrl.GRPCServer{},
		Config: ctrl.Config{
			LoadFromConfigCenter: true,
			RemoteConfigAddr:     "127.0.0.1:8500",
			RemoteConfigToken:    "73393977-9298-f39f-aaed-c3ab10e04e4d",
		},
		HealthProbeBindAddress: ":8083",
	})
	if err != nil {
		panic(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(err)
	}

	if err = mgr.Start(context.Background()); err != nil {
		panic(err)
	}
}

type Test struct {
	redisClient *redis.Client
	logger      *zap.Logger
	enc         encryption.Encryptor
}

func (t *Test) Init(cfg *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password, // no password set
		DB:       0,                  // use default DB
	})
	t.redisClient = rdb
	t.enc = encryption.NewEncryptor([]byte(cfg.Encryption.Passphrase), cfg.Encryption.Name, cfg.Encryption.Email, cfg.Encryption.RsaBits, cfg.Encryption.Enable)
	t.logger = log.NewDevLogger("test")
	return nil
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Version() string {
	return version.FullVersion()
}

func (t *Test) Registry(r gin.IRoutes) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func (t *Test) RegistryWithMiddle(r gin.IRoutes) {}
