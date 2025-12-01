package types

import (
	"net/http"
	"net/url"
	"sync"
	"time"
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
	Weight       int               `yaml:"weight" mapstructure:"weight"`
	
	//health check configs
	HealthCheck  HealthCheckConfig `yaml:"health_check" mapstructure:"health_check"`
}

type Server struct {
	URL            *url.URL   
	HealthCheckURL *url.URL
	Proxy          http.Handler 

	Weight         int
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

func (sc *ServerConfig) Validate() error {

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