package types

import (
	"net/url"
	"sync"
	"time"
)

// config struct for configuration
// includes LoadBalancer, Servers, Strategy
type Config struct {
	LoadBalancer LoadBalancer 	`yaml:"load_balancer"`
	Strategy string 	`yaml:"strategy"`
	ServerConfig []ServerConfig 	`yaml:"servers"`
}

type LoadBalancer struct {
		Current int 		`yaml:"current"`
		Mutex sync.Mutex 	`yaml:"mutex"`
	}


type Server struct {
	URL *url.URL `yaml:"url"`
	Threshold int `yaml:"threshold"`
	HealthEndpoint string `yaml:"healthEndpoint"`
	IsHealthy bool `yaml:"is_healthy"`
	Mutex sync.Mutex `yaml:"mutex"`
}

type ServerConfig struct {
	URL string `yaml:"url"`
	
	HealthCheck struct {
		Protocol string
		Port int
		Path string
		Timeout time.Duration
		Interval time.Duration
		HealthyThreshold int
		UnhealthyThreshold int
		Code int
	}
}