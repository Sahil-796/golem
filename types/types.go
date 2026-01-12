package types

import (
	"net/http"
	"net/url"
	"sync"
	"time"
	"fmt"
)

type Config struct {
	Strategy string         `yaml:"strategy" mapstructure:"strategy"` 
	Servers  []ServerConfig `yaml:"server_configs" mapstructure:"server_configs"`
}

type ServerConfig struct {
	// Server information
	Host         string            `yaml:"host" mapstructure:"host"`
	Port         int               `yaml:"port" mapstructure:"port"`
	Protocol     string            `yaml:"protocol" mapstructure:"protocol"` // http, https

	// Settings or configs
	ProxyTimeout time.Duration     `yaml:"proxy_timeout" mapstructure:"proxy_timeout"`
	Weight       int               `yaml:"weight" mapstructure:"weight"` // weight for weighted rr
	
	//health check configs
	HealthCheck  HealthCheckConfig `yaml:"health_check" mapstructure:"health_check"`
}

func (sc *ServerConfig) Validate() error {
	if sc.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if sc.Port <= 0 {
		return fmt.Errorf("server port must be a positive integer")
	}
	if sc.Protocol == "" {
		sc.Protocol = "http" // Default
	}
	if sc.Protocol != "http" && sc.Protocol != "https" {
		return fmt.Errorf("protocol must be 'http' or 'https'")
	}

	// Validate the nested health check config
	if err := sc.HealthCheck.Validate(); err != nil {
		return fmt.Errorf("health check validation error: %w", err)
	}
	return nil
}

type Server struct {
	URL            *url.URL   
	HealthCheckURL *url.URL
	Proxy          http.Handler  // handler interface is pre-set and is boxed as NewSingleHostReverseProxy()

	Weight         int
	CurrentWeight  int // current weight/laod on the server, starts from 0, reduces
	CurrentConnections int // current number of active connections on the server, starts from 0, increases
	IsHealthy      bool       
	Status         string
	ConsecutiveFailures  int 
	ConsecutiveSuccesses int 
	Mutex                sync.Mutex 
	LastCheck            time.Time
}

type HealthCheckConfig struct {

	Path               string        `yaml:"path" mapstructure:"path"` 
	Timeout            time.Duration `yaml:"timeout" mapstructure:"timeout"`
	Interval           time.Duration `yaml:"interval" mapstructure:"interval"`
	HealthyThreshold   int           `yaml:"healthy_threshold" mapstructure:"healthy_threshold"`
	UnhealthyThreshold int           `yaml:"unhealthy_threshold" mapstructure:"unhealthy_threshold"`
	Code               int           `yaml:"code" mapstructure:"code"`
}

func (hc *HealthCheckConfig) Validate() error {
	if hc.Path == "" {
		return fmt.Errorf("health_check.path is required")
	}
	if hc.HealthyThreshold <= 0 {
		return fmt.Errorf("healthy_threshold must be > 0")
	}
	if hc.UnhealthyThreshold <= 0 {
		return fmt.Errorf("unhealthy_threshold must be > 0")
	}
	return nil
}

func (hc *HealthCheckConfig) SetDefaults() {
	if hc.Path == "" {
		hc.Path = "/"
	}
	if hc.Timeout == 0 {
		hc.Timeout = 5 * time.Second
	}
	if hc.Interval == 0 {
		hc.Interval = 10 * time.Second
	}
	if hc.HealthyThreshold == 0 {
		hc.HealthyThreshold = 3
	}
	if hc.UnhealthyThreshold == 0 {
		hc.UnhealthyThreshold = 3
	}
	if hc.Code == 0 {
		hc.Code = 200
	}
}