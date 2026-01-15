package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
	"net/http"
	"log"
)

type WeightedRoundRobin struct {
	index int
	Mutex sync.Mutex
} 

func (wrr *WeightedRoundRobin) Next(r *http.Request, servers []*types.Server) *types.Server {
	
	if r == nil {
		log.Printf("[WARN] WeightedRoundRobin.Next: received nil request")
		return nil
	}

	if len(servers) == 0 {
		log.Printf("[WARN] WeightedRoundRobin.Next: no servers available")
		return nil
	}
	
	wrr.Mutex.Lock()
	defer wrr.Mutex.Unlock()
	n := len(servers)
	
	for range n {
		server := servers[wrr.index]
		
		if server == nil {
			log.Printf("[ERROR] WeightedRoundRobin.Next: server at index %d is nil", wrr.index)
			wrr.index = (wrr.index + 1) % n
			continue
		}
		
		server.Mutex.Lock()
		healthy := server.IsHealthy
		currentWeight := server.CurrentWeight
		weight := server.Weight
		
		if weight < 0 {
			log.Printf("[ERROR] WeightedRoundRobin.Next: server at index %d has negative weight %d", wrr.index, weight)
			server.Mutex.Unlock()
			wrr.index = (wrr.index + 1) % n
			continue
		}

		if weight == 0 {
			log.Printf("[WARN] WeightedRoundRobin.Next: server at index %d has zero weight, skipping", wrr.index)
			server.Mutex.Unlock()
			wrr.index = (wrr.index + 1) % n
			continue
		}
		
		if currentWeight >= weight {
			server.CurrentWeight = 0 
			server.Mutex.Unlock()
			wrr.index = (wrr.index + 1) % n
			continue
		}
		
		
		if !healthy {
			server.Mutex.Unlock()
			wrr.index = (wrr.index + 1) % n
			continue
		}
		
		server.CurrentWeight++
		server.Mutex.Unlock()
		
		log.Printf("[DEBUG] WeightedRoundRobin.Next: selected server at index %d (weight: %d/%d)", wrr.index, server.CurrentWeight, weight)
		return server	}
	
	log.Printf("[WARN] WeightedRoundRobin.Next: no suitable server found after checking %d servers",n)
	return nil
}
