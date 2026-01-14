package strategy

import (
	"net"
	"net/http"
	"strings"
	"errors"
)

func getIP(r *http.Request) (string, error) {

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