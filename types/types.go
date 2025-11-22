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

type LoadBalancer struct {
	Current int        
	Mutex   sync.Mutex 
	Backends []*Server
}


// holds user relevant data for each server
type Server struct {
	URL            *url.URL   
	IsHealthy      bool       
	Status         string
	ConsecutiveFailures int 
	ConsecutiveSuccesses int 
	Mutex          sync.Mutex 
}

// holds internally relevant or config data for each server
type ServerConfig struct {
	URL string `yaml:"url" mapstructure:"url"`

	HealthCheckConfig HealthCheckConfig `yaml:"health_check" mapstructure:"health_check"`
}

type HealthCheckConfig struct {
	Host               string        `yaml:"host" mapstructure:"host"` // required if URL missing
	Protocol           string        `yaml:"protocol" mapstructure:"protocol"` // default: http
	Port               int           `yaml:"port" mapstructure:"port"` // required if URL missing
	Path               string        `yaml:"path" mapstructure:"path"` // default: /
	Timeout            time.Duration `yaml:"timeout" mapstructure:"timeout"`
	Interval           time.Duration `yaml:"interval" mapstructure:"interval"`
	HealthyThreshold   int           `yaml:"healthy_threshold" mapstructure:"healthy_threshold"`
	UnhealthyThreshold int           `yaml:"unhealthy_threshold" mapstructure:"unhealthy_threshold"`
	Code               int           `yaml:"code" mapstructure:"code"`
}


func (sc *ServerConfig) Validate() error {
	hc := sc.HealthCheckConfig

	if sc.URL != "" {
		return nil 
	}

	if hc.Host == "" {
		return fmt.Errorf("either 'url' OR 'health_check.host' must be provided")
	}

	if hc.Protocol == "tcp" && hc.Port == 0 {
		return fmt.Errorf("tcp requires explicit port in health_check.port")
	}

	// 4: HTTP / HTTPS → port can be missing (Docker / service discovery case)
	//       BUT give a warning later (not here)
	//       This is valid, let defaults or buildHealthURL handle it

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
		hc.HealthyThreshold = 1
	}
	if hc.UnhealthyThreshold == 0 {
		hc.UnhealthyThreshold = 1
	}
	if hc.Code == 0 {
		hc.Code = 200
	}
}
