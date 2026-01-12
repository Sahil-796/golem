package config

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Sahil-796/golem/types"
	"github.com/spf13/viper"
)

func LoadConfig() (*types.Config, []*types.Server, error) {

	// setting up viper
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("error reading config: %w", err)
	}

	cfg := &types.Config{}

	// viper unarmshal: yaml should match exact type to be unmarshaled
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, nil, fmt.Errorf("error reading config.yaml: %w", err)
	}

	// initialising run time server list
	runtimeServers := make([]*types.Server, 0, len(cfg.Servers))

	for i := range cfg.Servers {

		serverConfig := &cfg.Servers[i] // server configurations

		// Validation of the server (Top level) (Information)
		if serverConfig.Host == "" {
			return nil, nil, fmt.Errorf("server host is required")
		}
		if serverConfig.Port == 0 {
			return nil, nil, fmt.Errorf("server port is required")
		}
		if serverConfig.Protocol == "" {
			serverConfig.Protocol = "http" // Default
		}

		// Validation of Health Check
		serverConfig.HealthCheck.SetDefaults()
		if err := serverConfig.Validate(); err != nil {
			return nil, nil, fmt.Errorf("invalid health check config: %w", err)
		}

		// Build Health URL
		// passing parent fields + (path = health path)
		healthURL, err := BuildURL(
			serverConfig.Protocol,
			serverConfig.Host,
			serverConfig.Port,
			serverConfig.HealthCheck.Path,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid health url: %w", err)
		}

		// Build Target URL (Proxy Traffic)
		// passing parent fields + root
		targetURL, err := BuildURL(
			serverConfig.Protocol,
			serverConfig.Host,
			serverConfig.Port,
			"",
		)

		// Create runtime Server
		server := &types.Server{ // Ptr to server
			URL:            targetURL,
			HealthCheckURL: healthURL,
			// Proxy set later
			Weight:               serverConfig.Weight,
			CurrentConnections:   0,
			IsHealthy:            true,
			Status:               "initial",
			ConsecutiveFailures:  0,
			ConsecutiveSuccesses: 0,
		}

		// Create Proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		timeout := serverConfig.ProxyTimeout
		if timeout == 0 {
			timeout = 10 * time.Second
		}

		proxy.Transport = &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: timeout,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}

		// proxy error handler function
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[%s] Proxy Error: %v", targetURL.Host, err)

			server.Mutex.Lock()
			server.ConsecutiveFailures++
			server.ConsecutiveSuccesses = 0
			// Circuit Breaker Logic
			if server.ConsecutiveFailures >= serverConfig.HealthCheck.UnhealthyThreshold {
				server.IsHealthy = false
				server.Status = "Unhealthy"
				log.Printf("[%s] Marked Unhealthy due to proxy failure", targetURL.Host)
			}
			server.Mutex.Unlock()

			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf(`{"error": "Backend Unavailable", "details": "%v"}`, err)))
		}

		// Assign proxy to server
		server.Proxy = proxy

		runtimeServers = append(runtimeServers, server)
	}

	return cfg, runtimeServers, nil
}

// A generic helper to build URLs
func BuildURL(protocol, host string, port int, path string) (*url.URL, error) {
	fullHost := fmt.Sprintf("%s:%d", host, port)

	built := &url.URL{
		Scheme: protocol,
		Host:   fullHost,
		Path:   path,
	}
	return built, nil
}
