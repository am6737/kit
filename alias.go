package kit

import "github.com/cossim/kit/pkg/manager"

// Manager initializes shared dependencies such as Caches and Clients, and provides them to Runnables.
// A Manager is required to create Controllers.
type Manager = manager.Manager

// Options are the arguments for creating a new Manager.
type Options = manager.Options

type HTTPServer = manager.HttpServer

type GRPCServer = manager.GrpcServer

type Config = manager.Config

var (
	// NewManager returns a new Manager for creating Controllers.
	NewManager = manager.New
)
