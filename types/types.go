package types

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

// config struct for configuration
// includes LoadBalancer, Servers, Strategy
type Config struct {
	Strategy string         `yaml:"strategy" mapstructure:"strategy"` 
	Servers  []ServerConfig `yaml:"server_configs" mapstructure:"server_configs"`
}


// holds user relevant data for each server
type Server struct {
	URL            *url.URL   
	HealthCheckURL *url.URL
	IsHealthy      bool       
	Status         string
	ConsecutiveFailures int 
	ConsecutiveSuccesses int 
	Mutex          sync.Mutex 
	LastCheck      time.Time
}

// holds internally relevant or config data for each server
type ServerConfig struct {
	HealthCheckConfig HealthCheckConfig `yaml:"health_check" mapstructure:"health_check"`
}

type HealthCheckConfig struct {
	Host               string        `yaml:"host" mapstructure:"host"` 
	Protocol           string        `yaml:"protocol" mapstructure:"protocol"` // default: http
	Port               int           `yaml:"port" mapstructure:"port"` 
	Path               string        `yaml:"path" mapstructure:"path"` // default: /
	Timeout            time.Duration `yaml:"timeout" mapstructure:"timeout"`
	Interval           time.Duration `yaml:"interval" mapstructure:"interval"`
	HealthyThreshold   int           `yaml:"healthy_threshold" mapstructure:"healthy_threshold"`
	UnhealthyThreshold int           `yaml:"unhealthy_threshold" mapstructure:"unhealthy_threshold"`
	Code               int           `yaml:"code" mapstructure:"code"`
}


func (sc *ServerConfig) Validate() error {
	hc := sc.HealthCheckConfig


	if hc.Host == "" {
		return fmt.Errorf("'health_check.host' must be provided")
	}
	
	validProtocols := map[string]bool{"http": true, "https": true, "tcp": false, "udp": false}
	if!validProtocols[hc.Protocol] {
		return fmt.Errorf("'health_check.protocol' must be one of %v", validProtocols)
	}
	
	if hc.HealthyThreshold <= 0 {
        return fmt.Errorf("'healthy_threshold' must be a positive integer (>= 1)")
    }
    if hc.UnhealthyThreshold <= 0 {
        return fmt.Errorf("'unhealthy_threshold' must be a positive integer (>= 1)")
    }

	return nil
}

func (hc *HealthCheckConfig) SetDefaults() {
	if hc.Protocol == "" {
		hc.Protocol = "http"
	}
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
