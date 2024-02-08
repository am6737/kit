package cluster

import (
	"context"
	"github.com/cossim/kit/pkg/config"
	"github.com/go-logr/logr"
	"go.etcd.io/etcd/client/v2"
	"google.golang.org/grpc"
	"net/http"
)

type cluster struct {
	// config is the rest.config used to talk to the apiserver.  Required.
	config *config.Config

	httpClient *http.Client
	grpcClient *grpc.ClientConn
	client     client.Client

	// mapper is used to map resources to kind, and map kind and version.
	//mapper meta.RESTMapper

	// Logger is the logger that should be used by this manager.
	// If none is set, it defaults to log.Log global logger.
	logger logr.Logger
}

func (c *cluster) GetGRPCClient() *grpc.ClientConn {
	return c.grpcClient
}

func (c *cluster) GetConfig() *config.Config {
	return c.config
}

func (c *cluster) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *cluster) GetClient() client.Client {
	return c.client
}

//func (c *cluster) GetRESTMapper() meta.RESTMapper {
//	return c.mapper
//}

func (c *cluster) GetLogger() logr.Logger {
	return c.logger
}

func (c *cluster) Start(ctx context.Context) error {
	//defer c.recorderProvider.Stop(ctx)
	//return c.cache.Start(ctx)
	return nil
}
