package strategy

import (
	"errors"
	"hash/fnv"
	"net"
	"net/http"
	"strings"
	"sync"
	"github.com/Sahil-796/golem/types"
	"log"
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
	
	if r == nil {
		return "", errors.New("request is nil")
	}
	
	// a client request mighhtt have xff if it has a proxy in middle

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		
		parts := strings.SplitSeq(xff, ",")
		for part := range parts {
			ip := strings.TrimSpace(part) // trim whitespace
			if net.ParseIP(ip) != nil {
				log.Printf("[DEBUG] getIP: extracted IP from X-Forwarded-For: %s", ip)
				return ip, nil
			} else {
				log.Printf("[WARN] getIP: invalid IP in X-Forwarded-For: %s", ip)
			}
		}
	}
	
	// fallback using RemoteAddr if browser connects to golem directly
	if r.RemoteAddr == "" {
		return "", errors.New("RemoteAddr is empty")
	}
	
	ip, _, err := net.SplitHostPort(r.RemoteAddr) 
		
	if err != nil {
		if parsedIP := net.ParseIP(r.RemoteAddr); parsedIP != nil {
					log.Printf("[DEBUG] getIP: extracted IP from RemoteAddr (no port): %s", r.RemoteAddr)
					return parsedIP.String(), nil
				}
				return "", errors.New("invalid RemoteAddr")
	}
		
	netIP := net.ParseIP(ip) // giga chad
	if netIP == nil {
		return "", errors.New("invalid IP")
	}
	
	log.Printf("[DEBUG] getIP: extracted IP from RemoteAddr: %s", netIP.String())
	return netIP.String(), nil

}

func (i *IPHash) Next(request *http.Request, servers []*types.Server) *types.Server {
	
	if request == nil {
		log.Printf("[ERROR] IPHash.Next: received nil request")
		return nil
	}
	
	if len(servers) == 0 {
		log.Printf("[WARN] IPHash.Next: no servers available")
		return nil	}
	
	i.Mutex.Lock()
	defer i.Mutex.Unlock()
	
	clientIP, err := getIP(request)
	if err != nil {
		log.Printf("[ERROR] IPHash.Next: failed to get client IP: %v", err)
		return nil
	}
	
	var maxScore uint64
	var selectedServer *types.Server
	
	for idx, server := range servers {
		
		if server == nil {
			log.Printf("[ERROR] IPHash.Next: server at index %d is nil", idx)
			continue
		}
		
		server.Mutex.Lock()
		isHealthy := server.IsHealthy
		serverURL := server.URL
		server.Mutex.Unlock()

		if serverURL == nil {
			log.Printf("[ERROR] IPHash.Next: server at index %d has nil URL", idx)
			continue
		}

		if !isHealthy {
			continue
		}

		score := hrwScore(clientIP, serverURL.String())
		
		if score > maxScore {
			maxScore = score
			selectedServer = server
		}

	}
	
	if selectedServer != nil {
		log.Printf("[DEBUG] IPHash.Next: selected server for client IP %s with score %d", clientIP, maxScore)
	} else {
		log.Printf("[WARN] IPHash.Next: no healthy server found for client IP %s among %d servers", clientIP, len(servers))
	}
	
	return selectedServer
}