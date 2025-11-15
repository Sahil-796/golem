package golem

servers := [...]Server{from config}

func (lb* LoadBalancer) getNextServer(server []*Server) *Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()
	
	for i:=0; i<len(servers); i++ {
		
	}
}