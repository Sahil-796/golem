package strategy

import (
	"errors"
	"hash/fnv"
	"net"
	"net/http"
	"strings"
	"sync"
	"github.com/Sahil-796/golem/types"
)

type IPHash struct {
	Mutex sync.Mutex
}

// using consistent hashing (Rendezvous hashing) over normal hashing
// this ensures consistent distribution of all ip
func hrwScore(clientIP, backendURL string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(clientIP))
	h.Write([]byte(backendURL))

	score := h.Sum64()


	return score
}


func getIP(r *http.Request) (string, error) {
	
	// a client request mighhtt have xff if it has a proxy in middle

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		
		parts := strings.SplitSeq(xff, ",")
		for part := range parts {
			ip := strings.TrimSpace(part) // trim whitespace
			if net.ParseIP(ip) != nil {
				return ip, nil
			}
		}
	}
	
	// fallback using RemoteAddr if browser connects to golem directly
	ip, _, err := net.SplitHostPort(r.RemoteAddr) 
		
	if err != nil {
		return "", err
	}
		
	netIP := net.ParseIP(ip) // giga chad
	if netIP == nil {
		return "", errors.New("invalid IP")
	}
	
	return netIP.String(), nil

}

func (i *IPHash) Next(request *http.Request, servers []*types.Server) *types.Server {
	if len(servers) == 0 {
		return nil
	}
	
	i.Mutex.Lock()
	defer i.Mutex.Unlock()
	
	clientIP, err := getIP(request)
	if err != nil {
		return nil
	}
	
	var maxScore uint64
	var selectedServer *types.Server
	
	for _, server := range servers {
		
		server.Mutex.Lock()
		score := hrwScore(clientIP, server.URL.String())
		server.Mutex.Unlock()
		
		if score > maxScore {
			maxScore = score
			selectedServer = server
		}

	}
	return selectedServer
}