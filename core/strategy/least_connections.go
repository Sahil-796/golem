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
			continue
		}
		
		server.Mutex.Lock()
		serverConnections := server.CurrentConnections;
		server.Mutex.Unlock()
		
		
		if currConnections < serverConnections {
			server = runtimeServer
		}
	}
	return server
}