package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
	"net/http"
	"log"
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
	defer lc.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	var server *types.Server
	var minConnections int
	
	for i := range n {
		
		runtimeServer := servers[i]
		
		if runtimeServer == nil {
			log.Printf("[ERROR] LeastConnections.Next: server at index %d is nil", i)
			continue
		}
		
		runtimeServer.Mutex.Lock()
		isHealthy := runtimeServer.IsHealthy
		currConnections := runtimeServer.CurrentConnections
		runtimeServer.Mutex.Unlock()
		
		if !isHealthy {
			continue
		}
		
		if currConnections < 0 {
			log.Printf("[ERROR] LeastConnections.Next: server at index %d has negative connection count %d", i, currConnections)
			continue
		}
		
		if server == nil {
			server = runtimeServer
			minConnections = currConnections
			continue
		}
		
		if currConnections < minConnections {
			server = runtimeServer
			minConnections = currConnections
		}
	}
	if server != nil {
		log.Printf("[DEBUG] LeastConnections.Next: selected server with %d connections", minConnections)
	} else {
		log.Printf("[WARN] LeastConnections.Next: no healthy servers found among %d servers", n)
	}

	return server
}