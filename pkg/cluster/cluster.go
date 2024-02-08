package cluster

import (
	"context"
	"errors"
	"github.com/cossim/kit/pkg/client"
	"github.com/cossim/kit/pkg/config"
	"github.com/cossim/kit/pkg/log"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

// Cluster provides various methods to interact with a cluster.
type Cluster interface {
	GetHTTPClient() *http.Client

	GetGRPCClient() *grpc.ClientConn

	// Start starts the cluster
	Start(ctx context.Context) error
}

// Options are the possible options that can be configured for a Cluster.
type Options struct {
	// Logger is the logger that should be used by this Cluster.
	// If none is set, it defaults to log.Log global logger.
	Logger logr.Logger

	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	// there will a 10 percent jitter between the SyncPeriod of all controllers
	// so that all controllers will not send list requests simultaneously.
	//
	// Deprecated: Use Cache.SyncPeriod instead.
	SyncPeriod *time.Duration

	// HTTPClient is the http client that will be used to create the default
	// Cache and Client. If not set the rest.HTTPClientFor function will be used
	// to create the http client.
	HTTPClient *http.Client

	GRPCClient *grpc.ClientConn

	// Client is the client.Options that will be used to create the default Client.
	// By default, the client will use the cache for reads and direct calls for writes.
	Client client.Options
}

// Option can be used to manipulate Options.
type Option func(*Options)

// New constructs a brand new cluster.
func New(config *config.Config, opts ...Option) (Cluster, error) {
	if config == nil {
		return nil, errors.New("must specify Config")
	}

	originalConfig := config

	//config = rest.CopyConfig(config)
	//if config.UserAgent == "" {
	//	config.UserAgent = rest.DefaultKubernetesUserAgent()
	//}

	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}
	options, err := setOptionsDefaults(options, config)
	if err != nil {
		options.Logger.Error(err, "Failed to set defaults")
		return nil, err
	}

	// Create the client, and default its options.
	clientOpts := options.Client
	{
		if clientOpts.HTTPClient == nil {
			clientOpts.HTTPClient = options.HTTPClient
		}
		if clientOpts.GRPCClient == nil {
			clientOpts.GRPCClient = options.GRPCClient
		}
	}

	return &cluster{
		config:     originalConfig,
		httpClient: options.HTTPClient,
		//client:     clientWriter,
		logger: options.Logger,
	}, nil
}

// setOptionsDefaults set default values for Options fields.
func setOptionsDefaults(options Options, config *config.Config) (Options, error) {
	if options.HTTPClient == nil {
		var err error
		options.HTTPClient, err = config.HTTPClientFor(config)
		if err != nil {
			return options, err
		}
	}

	if options.GRPCClient == nil {
		var err error
		options.GRPCClient, err = config.GRPCClientFor(config)
		if err != nil {
			return options, err
		}
	}

	if options.Logger.GetSink() == nil {
		options.Logger = log.NewLogrLogger("cluster")
	}

	return options, nil
}
