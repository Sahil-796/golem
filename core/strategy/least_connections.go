package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
)

type LeastConnections struct {
	Mutex sync.Mutex
} 

func (lc *LeastConnections) Next(servers []*types.Server) *types.Server {
	if len(servers) == 0 {
		return nil
	} 
	
	lc.Mutex.Lock()
	defer lc.Mutex.Unlock() // defer to the end
	n:=len(servers)
	
	var server *types.Server
	var minConnections int
	
	for i := range n {
		
		runtimeServer := servers[i]
		runtimeServer.Mutex.Lock()
		isHealthy := runtimeServer.IsHealthy
		currConnections := runtimeServer.CurrentConnections
		runtimeServer.Mutex.Unlock()
		
		if !isHealthy {
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
	return server
}