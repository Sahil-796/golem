package main

import (
	"log"
	"net/http"

	"github.com/Sahil-796/golem/config"
	"github.com/Sahil-796/golem/core"
	"github.com/Sahil-796/golem/core/health"
	// "github.com/gin-gonic/gin"
	// "github.com/Sahil-796/golem/server/pkg/balancer"
	// "fmt"
)

func main() {
	
	cfg, servers, err := config.LoadConfig()
	
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	
	go health.StartHealthCheckers(servers, cfg.Servers)
	
	lb := core.NewLoadBalancer(cfg.Strategy, servers)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend := lb.Balance()
		
		if backend == nil {
			// status = 503 -> service unavailable
			http.Error(w, "No healthy backend available", http.StatusServiceUnavailable)
		}
		
		backend.Proxy.ServeHTTP(w, r) 
	})
	
	log.Println("Load balancer running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}