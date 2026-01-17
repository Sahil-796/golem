package strategy

import (
	"log"
	"net/http"
	"sync"

	"github.com/Sahil-796/golem/types"
)

type LeastConnections struct {
	Mutex sync.Mutex
}

func (lc *LeastConnections) Next(r *http.Request, servers []*types.Server) *types.Server {
	if r == nil {
		log.Printf("[WARN] LeastConnections.Next: received nil request")
		return nil
	}

	if len(servers) == 0 {
		log.Printf("[WARN] LeastConnections.Next: no servers available")
		return nil
	}

	lc.Mutex.Lock()
	defer lc.Mutex.Unlock()

	var chosen *types.Server
	min := int(^uint(0) >> 1) // max int

	for i := range servers {
		s := servers[i]
		if s == nil {
			continue
		}

		s.Mutex.Lock()
		healthy := s.IsHealthy
		conns := s.CurrentConnections
		s.Mutex.Unlock()

		if !healthy {
			continue
		}

		if conns < min {
			min = conns
			chosen = s
		}
	}

	if chosen != nil {
		log.Printf("[DEBUG] LeastConnections.Next: selected server with %d active connections", min)
	} else {
		log.Printf("[WARN] LeastConnections.Next: no healthy servers found")
	}

	return chosen
}
