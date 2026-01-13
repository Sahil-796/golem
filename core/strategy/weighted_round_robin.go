package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
)

type WeightedRoundRobin struct {
	index int
	Mutex sync.Mutex
} 

func (wrr *WeightedRoundRobin) Next(servers []*types.Server) *types.Server {
	
	if len(servers) == 0 {
		return nil
	} 
	
	wrr.Mutex.Lock()
	defer wrr.Mutex.Unlock()
	n := len(servers)
	
	for range n {
		server := servers[wrr.index]
		
		server.Mutex.Lock()
		healthy := server.IsHealthy
		currentWeight := server.CurrentWeight
		weight := server.Weight
		
		
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
		
		return server
	}
	
	return nil
}
