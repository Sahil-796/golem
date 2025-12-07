package health

import (
	"context"
	"log"
	"net/http"

	"github.com/Sahil-796/golem/types"
)

const (
    StatusWarmingUp = "Warming-Up"   // before thresholds reached
    StatusDegraded  = "Degraded"     // failing but not yet unhealthy
    StatusHealthy   = "Healthy"      // well congrats
    StatusUnhealthy = "Unhealthy"    // doomed
)

func ActiveCheckSingle(server *types.Server, ServerConfig types.ServerConfig)  {
	
	// server: the live list of server (mutex) which keeps track of health status of each server and is parallelized
	// ServerConfig: the configuration for each server's health check
	
	
	client := &http.Client{}
		
		healthUrl := server.HealthCheckURL

		ctx, cancel := context.WithTimeout(context.Background(), ServerConfig.HealthCheck.Timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodHead, healthUrl.String(), nil)
		if err != nil {
			log.Printf("[HealthCheck] %s request error: %v", healthUrl, err)
			return
		}

		resp, err := client.Do(req)

		server.Mutex.Lock()
		defer server.Mutex.Unlock()
		
		if err != nil || resp.StatusCode != http.StatusOK {
			
			log.Printf("[HealthCheck] %s failed: %v", server.HealthCheckURL, err)
			server.ConsecutiveFailures++
			server.ConsecutiveSuccesses = 0
			
			if server.ConsecutiveFailures >= ServerConfig.HealthCheck.UnhealthyThreshold {
					if server.IsHealthy {
							log.Printf("[HealthCheck] %s is now %s (Error: %v)", server.HealthCheckURL, StatusUnhealthy, err)
						}
						
						server.IsHealthy = false
						server.Status = StatusUnhealthy
					} else {
						server.Status = StatusDegraded
						log.Printf("[HealthCheck] %s is %s (%d/%d failures)", 
							server.HealthCheckURL, 
							StatusDegraded, 
							server.ConsecutiveFailures, 
							ServerConfig.HealthCheck.UnhealthyThreshold,
						)					
					}
			
			return 
		}
				
		server.ConsecutiveSuccesses++
		server.ConsecutiveFailures = 0
		
		if server.ConsecutiveSuccesses >= ServerConfig.HealthCheck.HealthyThreshold {
			if !server.IsHealthy {
				log.Printf("[HealthCheck] %s is now %s", server.HealthCheckURL, StatusHealthy)
			}
			
			server.IsHealthy = true
			server.Status = StatusHealthy

		} else {
			server.Status = StatusWarmingUp
			log.Printf("[HealthCheck] %s is %s (%d/%d successes)", 
				server.HealthCheckURL, 
				StatusWarmingUp, 
				server.ConsecutiveSuccesses, 
				ServerConfig.HealthCheck.HealthyThreshold,
			)		}

		
	
	
}

