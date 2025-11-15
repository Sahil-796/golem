package golem
import (
	"sync"
	"net/url"
)

type LoadBalancer struct {
	Current int
	Mutex sync.Mutex
}

type Server struct {
	URL *url.URL
	IsHealthy bool
	Mutex sync.Mutex
}