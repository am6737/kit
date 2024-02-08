package transport

import (
	"fmt"
	"net/http"
)

// WrapperFunc wraps an http.RoundTripper when a new transport
// is created for a client, allowing per connection behavior
// to be injected.
type WrapperFunc func(rt http.RoundTripper) http.RoundTripper

// Wrappers accepts any number of wrappers and returns a wrapper
// function that is the equivalent of calling each of them in order. Nil
// values are ignored, which makes this function convenient for incrementally
// wrapping a function.
func Wrappers(fns ...WrapperFunc) WrapperFunc {
	if len(fns) == 0 {
		return nil
	}
	// optimize the common case of wrapping a possibly nil transport wrapper
	// with an additional wrapper
	if len(fns) == 2 && fns[0] == nil {
		return fns[1]
	}
	return func(rt http.RoundTripper) http.RoundTripper {
		base := rt
		for _, fn := range fns {
			if fn != nil {
				base = fn(base)
			}
		}
		return base
	}
}

// New returns an http.RoundTripper that will provide the authentication
// or transport level security defined by the provided Config.
func New(config *Config) (http.RoundTripper, error) {
	// Set transport level security
	if config.Transport != nil && (config.HasCA() || config.HasCertAuth() || config.HasCertCallback() || config.TLS.Insecure) {
		return nil, fmt.Errorf("using a custom transport with TLS certificate options or the insecure flag is not allowed")
	}

	if !isValidHolders(config) {
		return nil, fmt.Errorf("misconfigured holder for dialer or cert callback")
	}

	var (
		rt http.RoundTripper
		//err error
	)

	if config.Transport != nil {
		rt = config.Transport
	} else {
		//rt, err = tlsCache.get(config)
		//if err != nil {
		//	return nil, err
		//}
	}

	return HTTPWrappersForConfig(config, rt)
}

func isValidHolders(config *Config) bool {
	if config.TLS.GetCertHolder != nil && config.TLS.GetCertHolder.GetCert == nil {
		return false
	}

	if config.DialHolder != nil && config.DialHolder.Dial == nil {
		return false
	}

	return true
}

// HTTPWrappersForConfig wraps a round tripper with any relevant layered
// behavior from the config. Exposed to allow more clients that need HTTP-like
// behavior but then must hijack the underlying connection (like WebSocket or
// HTTP2 clients). Pure HTTP clients should use the RoundTripper returned from
// New.
func HTTPWrappersForConfig(config *Config, rt http.RoundTripper) (http.RoundTripper, error) {
	//if config.WrapTransport != nil {
	//	rt = config.WrapTransport(rt)
	//}
	//
	//rt = DebugWrappers(rt)
	//
	//// Set authentication wrappers
	//switch {
	//case config.HasBasicAuth() && config.HasTokenAuth():
	//	return nil, fmt.Errorf("username/password or bearer token may be set, but not both")
	//case config.HasTokenAuth():
	//	var err error
	//	rt, err = NewBearerAuthWithRefreshRoundTripper(config.BearerToken, config.BearerTokenFile, rt)
	//	if err != nil {
	//		return nil, err
	//	}
	//case config.HasBasicAuth():
	//	rt = NewBasicAuthRoundTripper(config.Username, config.Password, rt)
	//}
	//if len(config.UserAgent) > 0 {
	//	rt = NewUserAgentRoundTripper(config.UserAgent, rt)
	//}
	//if len(config.Impersonate.UserName) > 0 ||
	//	len(config.Impersonate.UID) > 0 ||
	//	len(config.Impersonate.Groups) > 0 ||
	//	len(config.Impersonate.Extra) > 0 {
	//	rt = NewImpersonatingRoundTripper(config.Impersonate, rt)
	//}
	return rt, nil
}
