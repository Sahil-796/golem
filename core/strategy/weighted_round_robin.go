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
		
		if server.CurrentWeight >= server.Weight {
			server.CurrentWeight = 0 
			server.Mutex.Unlock()
			wrr.index = (wrr.index + 1) % n
			continue
		}
		
		server.CurrentWeight++
		server.Mutex.Unlock()
		
		if healthy {
			return server
		}
		
		// If not healthy, continue to next server
		wrr.index = (wrr.index + 1) % n
	}
	
	return nil
}
