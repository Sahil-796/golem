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
	
	http.HandleFunc("/", func(writter http.ResponseWriter, request *http.Request) {
		backend := lb.Balance(request)
		
		if backend == nil {
			// status = 503 -> service unavailable
			http.Error(writter, "No healthy backend available", http.StatusServiceUnavailable)
		}
		
		backend.Proxy.ServeHTTP(writter, request) 
	})
	
	log.Println("Load balancer running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}