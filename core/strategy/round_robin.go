package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
	"net/http"
	"log"
)

type RoundRobin struct {
	index int
	Mutex sync.Mutex
} 

func (rr *RoundRobin) Next(r *http.Request, servers []*types.Server) *types.Server {
	
	if r == nil {
		log.Printf("[WARN] RoundRobin.Next: received nil request")
		return nil
	}

	if len(servers) == 0 {
		log.Printf("[WARN] RoundRobin.Next: no servers available")
		return nil
	}
	
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	for range n {
		
		rr.index = (rr.index +1) % n
		s := servers[rr.index]
		
		if s == nil {
			log.Printf("[ERROR] RoundRobin.Next: server at index %d is nil", rr.index)
			continue
		}
		
		s.Mutex.Lock()
		healthy := s.IsHealthy
		s.Mutex.Unlock()
		
		if healthy { 
			log.Printf("[DEBUG] RoundRobin.Next: selected server at index %d", rr.index)
			return s		
		}
	}
	log.Printf("[WARN] RoundRobin.Next: all %d servers are unhealthy", n)
	return nil
}