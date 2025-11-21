package types

import (
	"net/url"
	"sync"
	"time"
)

// config struct for configuration
// includes LoadBalancer, Servers, Strategy
type Config struct {
	Strategy     string         `yaml:"strategy"`
	Servers      []ServerConfig `yaml:"server_configs"`  //array for all server configs
}

type LoadBalancer struct {
	Current int        `yaml:"current"`
	Mutex   sync.Mutex `yaml:"mutex"`
	Backends []*Server `yaml:"backends"`
}


// holds user relevant data for each server
type Server struct {
	URL            *url.URL   `yaml:"url"`
	IsHealthy      bool       `yaml:"is_healthy"`
	Status         string     `yaml:"status"`
	ConsecutiveFailures int `yaml:"consecutive_failures"`
	ConsecutiveSuccesses int `yaml:"consecutive_successes"`
	Mutex          sync.Mutex `yaml:"mutex"`
}

// holds internally relevant or config data for each server
type ServerConfig struct {
	URL string `yaml:"url"`

	HealthCheckConfig HealthCheckConfig `yaml:"health_check"`
}

type HealthCheckConfig struct {
	Host 			   string 		 `yaml:"host"`
	Protocol           string        `yaml:"protocol"`
	Port               int           `yaml:"port"`
	Path               string        `yaml:"path"`
	Timeout            time.Duration `yaml:"timeout"`
	Interval           time.Duration `yaml:"interval"`
	HealthyThreshold   int           `yaml:"healthy_threshold"`
	UnhealthyThreshold int           `yaml:"unhealthy_threshold"`
	Code               int           `yaml:"code"`
}
