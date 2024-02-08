package metrics

import "context"

// Server is a server that serves metrics.
type Server interface {
	// NeedLeaderElection implements the LeaderElectionRunnable interface, which indicates
	// the metrics server doesn't need leader election.
	NeedLeaderElection() bool

	// Start runs the server.
	// It will install the metrics related resources depending on the server configuration.
	Start(ctx context.Context) error
}
