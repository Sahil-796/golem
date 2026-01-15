package core

import (
	"log"
	"net/http"
	"sync"

	"github.com/Sahil-796/golem/core/strategy"
	"github.com/Sahil-796/golem/types"
)

type LoadBalancer struct {
	Mutex    sync.Mutex
	Backends []*types.Server
	Strategy strategy.Strategy
}

func NewLoadBalancer(strategyName string, backends []*types.Server) *LoadBalancer {
	return &LoadBalancer{
		Strategy: strategy.Get(strategyName),
		Backends: backends,
	}
}

func (lb *LoadBalancer) Balance(r *http.Request) *types.Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	return lb.Strategy.Next(r, lb.Backends)
}

// ServeHTTP handles the request and manages connection counting
// increments the connection count before proxying and decrements after completion
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.Balance(r)

	if backend == nil {
		log.Printf("[ERROR] No healthy backend available for request: %s %s", r.Method, r.URL.Path)
		http.Error(w, "No healthy backend available", http.StatusServiceUnavailable)
		return
	}

	backend.Mutex.Lock()
	backend.CurrentConnections++
	currentConns := backend.CurrentConnections
	backend.Mutex.Unlock()

	log.Printf("[DEBUG] Proxying request to backend (active connections: %d)", currentConns)

	defer func() {
		backend.Mutex.Lock()
		backend.CurrentConnections--
		remainingConns := backend.CurrentConnections
		backend.Mutex.Unlock()

		log.Printf("[DEBUG] Request completed (remaining connections: %d)", remainingConns)
	}()

	// main proxy func
	backend.Proxy.ServeHTTP(w, r)
}
