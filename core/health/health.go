package core

import (
	"log"
	"net/http"
	"time"

	"github.com/Sahil-796/golem/types"
)

const (
    StatusWarmingUp = "Warming-Up"   // before thresholds reached
    StatusDegraded  = "Degraded"     // failing but not yet unhealthy
    StatusHealthy   = "Healthy"      // well congrats
    StatusUnhealthy = "Unhealthy"    // doomed
    StatusOffline   = "Offline"      // server is down/not reachable
)

func ActiveCheck(servers []*types.Server, ServerConfig []types.ServerConfig) error {
	
	// servers: the live list of servers (mutex) which keeps track of health status of each server and is parallelized
	// ServerConfig: the configuration for each server's health check
	
	for i := range servers {
		
		healthUrl := servers[i].HealthCheckURL
		
		// setting up net call for /health route
	
		client := &http.Client{
			Timeout: ServerConfig[i].HealthCheckConfig.Timeout,
		}

		result, err := client.Get(healthUrl.String())
		servers[i].LastCheck = time.Now()
		
		if result != nil {
			defer result.Body.Close()
		}
		
		servers[i].Mutex.Lock()
		if err != nil || result.StatusCode != http.StatusOK {
			
			log.Printf("[HealthCheck] %s failed: %v", servers[i].HealthCheckURL, err)
			servers[i].ConsecutiveFailures++
			servers[i].ConsecutiveSuccesses = 0
			
			if servers[i].ConsecutiveFailures >= ServerConfig[i].HealthCheckConfig.UnhealthyThreshold {
							servers[i].IsHealthy = false
							servers[i].Status = StatusUnhealthy
						} else {
							servers[i].Status = StatusDegraded
						}
			
			servers[i].Mutex.Unlock()
			continue  // don't return
		}
		
		servers[i].ConsecutiveSuccesses++
		servers[i].ConsecutiveFailures = 0
		
		if servers[i].ConsecutiveSuccesses >= ServerConfig[i].HealthCheckConfig.HealthyThreshold {
			servers[i].IsHealthy = true
			servers[i].Status = StatusHealthy
		} else {
			servers[i].Status = StatusWarmingUp
		}

		servers[i].Mutex.Unlock()

	}
	
	return nil 
}

