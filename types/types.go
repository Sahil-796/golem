package types
import (
	"sync"
	"net/url"
)
// config struct for configuration
// includes LoadBalancer, Servers, Strategy
type Config struct {
	LoadBalancer struct {
		Current int 		`yaml:"current"`
		Mutex sync.Mutex 	`yaml:"mutex"`
	} 					`yaml:"load_balancer"`
	Strategy string 	`yaml:"strategy"`
	Servers []Server 	`yaml:"servers"`
}


type Server struct {
	URL *url.URL `yaml:"url"`
	IsHealthy bool `yaml:"is_healthy"`
	Mutex sync.Mutex `yaml:"mutex"`
}