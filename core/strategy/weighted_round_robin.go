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
	defer wrr.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	for range n {
		
		server := servers[wrr.index]
		
		if (server.Weight <= server.CurrentWeight) {
			wrr.index = (wrr.index + 1) % n
			server = servers[wrr.index]

		}
		
		server.CurrentWeight++
		server.Mutex.Lock()
		healthy := server.IsHealthy
		server.Mutex.Unlock()
		
		if healthy { 
			return server
		}
		
	}
	
	return nil
}
