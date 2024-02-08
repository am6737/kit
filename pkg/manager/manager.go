package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/cossim/kit/pkg/cluster"
	"github.com/cossim/kit/pkg/config"
	"github.com/cossim/kit/pkg/discovery"
	"github.com/cossim/kit/pkg/healthz"
	"github.com/cossim/kit/pkg/server"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"net"
	"net/http"
)

// Manager initializes shared dependencies such as Caches and Clients, and provides them to Runnables.
// A Manager is required to create Controllers.
type Manager interface {
	// Cluster holds a variety of methods to interact with a cluster.
	cluster.Cluster

	// Add will set requested dependencies on the component, and cause the component to be
	// started when Start is called.  Add will inject any dependencies for which the argument
	// implements the inject interface - e.g. inject.Client.
	// Depending on if a Runnable implements LeaderElectionRunnable interface, a Runnable can be run in either
	// non-leaderelection mode (always running) or leader election mode (managed by leader election if enabled).
	Add(Runnable) error

	// Elected is closed when this manager is elected leader of a group of
	// managers, either because it won a leader election or because no leader
	// election was configured.
	Elected() <-chan struct{}

	// AddMetricsExtraHandler adds an extra handler served on path to the http httpServer that serves metrics.
	// Might be useful to register some diagnostic endpoints e.g. pprof. Note that these endpoints meant to be
	// sensitive and shouldn't be exposed publicly.
	// If the simple path -> handler mapping offered here is not enough, a new http httpServer/listener should be added as
	// Runnable to the manager via Add method.
	AddMetricsExtraHandler(path string, handler http.Handler) error

	// AddHealthzCheck allows you to add Healthz checker
	AddHealthzCheck(name string, check healthz.Checker) error

	// AddReadyzCheck allows you to add Readyz checker
	AddReadyzCheck(name string, check healthz.Checker) error

	// Start starts all registered Controllers and blocks until the context is cancelled.
	// Returns an error if there is an error starting any controller.
	//
	// If LeaderElection is used, the binary must be exited immediately after this returns,
	// otherwise components that need leader election might continue to run after the leader
	// lock was lost.
	Start(ctx context.Context) error

	// GetLogger returns this manager's logger.
	GetLogger() logr.Logger
}

type GrpcServer struct {
	MetricsBindAddress  string
	HealthzCheckAddress string
	server.GRPCService
}

type HttpServer struct {
	MetricsBindAddress  string
	HealthzCheckAddress string
	server.HTTPService
}

type Config struct {
	LocalDir             string         // 从本地路径加载配置
	LoadFromConfigCenter bool           // 通过配置中心获取配置
	RemoteConfigAddr     string         // 配置中心地址
	RemoteConfigToken    string         // 访问token
	Watcher              WatcherManager // 配置文件监控
}

type WatcherManager struct {
	Enabled bool
	DisCC   discovery.ConfigCenter
	Keys    []string // 监控的文件

	UpdateCh chan struct {
		Key   string
		Value interface{}
	}
}

type Options struct {
	Grpc   GrpcServer
	Http   HttpServer
	Logger zap.Logger
	Config Config

	readinessEndpointName string

	// Liveness probe endpoint name
	livenessEndpointName string

	// HealthProbeBindAddress is the TCP address that the controller should bind to
	// for serving health probes
	// It can be set to "0" or "" to disable serving the health probe.
	HealthProbeBindAddress string

	newHealthProbeListener func(addr string) (net.Listener, error)
}

// New returns a new Manager for creating Controllers.
// Note that if ContentType in the given config is not set, "application/vnd.kubernetes.protobuf"
// will be used for all built-in resources of Kubernetes, and "application/json" is for other types
// including all CRD resources.
func New(cfg *config.Config, opts Options) (Manager, error) {
	//if config == nil {
	//	return nil, errors.New("must specify Config")
	//}

	// Set default values for options fields
	opts = setOptionsDefaults(opts)

	newCfg := &config.Config{}
	if opts.Config.LoadFromConfigCenter {
		opts.Config.Watcher = loadRemoteConfig(opts.Config, newCfg)
	}

	//hs := server.NewHttpService(newCfg, opts.Http)

	// Create health probes listener. This will throw an error if the bind
	// address is invalid or already in use.
	healthProbeListener, err := opts.newHealthProbeListener(opts.HealthProbeBindAddress)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	runnables := newRunnables(context.Background, errChan)

	fmt.Println("sss")

	return &controllerManager{
		runnables: runnables,
		errChan:   errChan,
		//httpServer:              hs,
		readinessEndpointName:   opts.readinessEndpointName,
		livenessEndpointName:    opts.livenessEndpointName,
		httpHealthProbeListener: healthProbeListener,
	}, nil
}

func setOptionsDefaults(opts Options) Options {
	if opts.readinessEndpointName == "" {
		opts.readinessEndpointName = defaultReadinessEndpoint
	}

	if opts.livenessEndpointName == "" {
		opts.livenessEndpointName = defaultLivenessEndpoint
	}

	if opts.newHealthProbeListener == nil {
		opts.newHealthProbeListener = defaultHealthProbeListener
	}

	return opts
}

// defaultHealthProbeListener creates the default health probes listener bound to the given address.
func defaultHealthProbeListener(addr string) (net.Listener, error) {
	if addr == "" || addr == "0" {
		return nil, nil
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error listening on %s: %w", addr, err)
	}
	return ln, nil
}

func loadRemoteConfig(cfg Config, ac *config.Config) WatcherManager {
	cc, err := discovery.NewDefaultRemoteConfigManager(cfg.RemoteConfigAddr, discovery.WithToken(cfg.RemoteConfigToken))
	if err != nil {
		panic(err)
	}
	c, err := cc.Get("interface/user")
	if err != nil {
		panic(err)
	}
	*ac = *c

	return WatcherManager{
		Enabled:  false,
		Keys:     nil,
		UpdateCh: nil,
	}
}

// BaseContextFunc is a function used to provide a base Context to Runnables
// managed by a Manager.
type BaseContextFunc func() context.Context

// Runnable allows a component to be started.
// It's very important that Start blocks until
// it's done running.
type Runnable interface {
	// Start starts running the component.  The component will stop running
	// when the context is closed. Start blocks until the context is closed or
	// an error occurs.
	Start(context.Context) error
}

// runnables handles all the runnables for a manager by grouping them accordingly to their
// type (webhooks, caches etc.).
type runnables struct {
	HTTPServers *runnableGroup
	GRPCServers *runnableGroup
}

// StopAndWait waits for all the runnables to finish before returning.
func (r *runnableGroup) StopAndWait(ctx context.Context) {
	r.stopOnce.Do(func() {
		// Close the reconciler channel once we're done.
		defer close(r.ch)

		_ = r.Start(ctx)
		r.stop.Lock()
		// Store the stopped variable so we don't accept any new
		// runnables for the time being.
		r.stopped = true
		r.stop.Unlock()

		// Cancel the internal channel.
		r.cancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			// Wait for all the runnables to finish.
			r.wg.Wait()
		}()

		select {
		case <-done:
			// We're done, exit.
		case <-ctx.Done():
			// Calling context has expired, exit.
		}
	})
}

// Start starts the group and waits for all
// initially registered runnables to start.
// It can only be called once, subsequent calls have no effect.
func (r *runnableGroup) Start(ctx context.Context) error {
	var retErr error

	r.startOnce.Do(func() {
		defer close(r.startReadyCh)

		// Start the internal reconciler.
		go r.reconcile()

		// Start the group and queue up all
		// the runnables that were added prior.
		r.start.Lock()
		r.started = true
		for _, rn := range r.startQueue {
			rn.signalReady = true
			r.ch <- rn
		}
		r.start.Unlock()

		// If we don't have any queue, return.
		if len(r.startQueue) == 0 {
			return
		}

		// Wait for all runnables to signal.
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); !errors.Is(err, context.Canceled) {
					retErr = err
				}
			case rn := <-r.startReadyCh:
				for i, existing := range r.startQueue {
					if existing == rn {
						// Remove the item from the start queue.
						r.startQueue = append(r.startQueue[:i], r.startQueue[i+1:]...)
						break
					}
				}
				// We're done waiting if the queue is empty, return.
				if len(r.startQueue) == 0 {
					return
				}
			}
		}
	})

	return retErr
}

// reconcile is our main entrypoint for every runnable added
// to this group. Its primary job is to read off the internal channel
// and schedule runnables while tracking their state.
func (r *runnableGroup) reconcile() {
	for runnable := range r.ch {
		fmt.Println("runnable => ", runnable)
		// Handle stop.
		// If the shutdown has been called we want to avoid
		// adding new goroutines to the WaitGroup because Wait()
		// panics if Add() is called after it.
		{
			r.stop.RLock()
			if r.stopped {
				// Drop any runnables if we're stopped.
				r.errChan <- errRunnableGroupStopped
				r.stop.RUnlock()
				continue
			}

			// Why is this here?
			// When StopAndWait is called, if a runnable is in the process
			// of being added, we could end up in a situation where
			// the WaitGroup is incremented while StopAndWait has called Wait(),
			// which would result in a panic.
			r.wg.Add(1)
			r.stop.RUnlock()
		}

		// Start the runnable.
		go func(rn *readyRunnable) {
			go func() {
				if rn.Check(r.ctx) {
					if rn.signalReady {
						r.startReadyCh <- rn
					}
				}
			}()

			// If we return, the runnable ended cleanly
			// or returned an error to the channel.
			//
			// We should always decrement the WaitGroup here.
			defer r.wg.Done()

			// Start the runnable.
			if err := rn.Start(r.ctx); err != nil {
				r.errChan <- err
			}
		}(runnable)
	}
}

var (
	errRunnableGroupStopped = errors.New("can't accept new runnable as stop procedure is already engaged")
)
