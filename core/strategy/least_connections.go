package strategy

import (
	"github.com/Sahil-796/golem/types"
	"sync"
)

type LeastConnections struct {
	index int
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
		
		if (server!=nil && servers[i].CurrentConnections > server.CurrentConnections) {
			// if (server == nil) {
				
			// }
			// server = servers[i]
		}
		
		
	}
	return server
}