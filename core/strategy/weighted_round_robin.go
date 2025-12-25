package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
)

type WeightedRoundRobin struct {
	index int
	Mutex sync.Mutex
} 

func (rr *WeightedRoundRobin) Next(servers []*types.Server) *types.Server {
	
	if len(servers) == 0 {
		return nil
	} 
	
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	for range n {
		
		
	}
	
	return nil
}