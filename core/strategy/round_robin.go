package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
)

type RoundRobin struct {
	index int
	Mutex sync.Mutex
} 

func (rr *RoundRobin) NewRoundRobin () *RoundRobin {
	return new(RoundRobin)
}

func (rr *RoundRobin) Next(servers []*types.Server) *types.Server {
	
	if len(servers) == 0 {
		return nil
	} 
	
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	for range(n) {
		
		s := servers[rr.index]
		rr.index = (rr.index +1) % n
		
		s.Mutex.Lock()
		healthy := s.IsHealthy
		s.Mutex.Unlock()
		
		if healthy { 
			return s
		}
	}
	
	return nil
}