package health

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

func ActiveCheckSingle(servers *types.Server, ServerConfig types.ServerConfig)  {
	
	// servers: the live list of servers (mutex) which keeps track of health status of each server and is parallelized
	// ServerConfig: the configuration for each server's health check
	
	
	client := &http.Client{
		Timeout: ServerConfig.HealthCheckConfig.Timeout,
	}
		
		healthUrl := servers.HealthCheckURL

		result, err := client.Get(healthUrl.String())
		servers.LastCheck = time.Now()
		
		if result != nil {
			defer result.Body.Close()
		}
		
		servers.Mutex.Lock()
		if err != nil || result.StatusCode != http.StatusOK {
			
			log.Printf("[HealthCheck] %s failed: %v", servers.HealthCheckURL, err)
			servers.ConsecutiveFailures++
			servers.ConsecutiveSuccesses = 0
			
			if servers.ConsecutiveFailures >= ServerConfig.HealthCheckConfig.UnhealthyThreshold {
							servers.IsHealthy = false
							servers.Status = StatusUnhealthy
							log.Printf("[HealthCheck] %s status set to unhealthy: %v", servers.HealthCheckURL, err)
						} else {
							servers.Status = StatusDegraded
						}
			
			servers.Mutex.Unlock()
			return 
		}
		
		log.Printf("[HealthCheck] %s Success: %v", servers.HealthCheckURL, err)
		
		servers.ConsecutiveSuccesses++
		servers.ConsecutiveFailures = 0
		
		if servers.ConsecutiveSuccesses >= ServerConfig.HealthCheckConfig.HealthyThreshold {
			servers.IsHealthy = true
			servers.Status = StatusHealthy
			log.Printf("[HealthCheck] %s status set to healthy: %v", servers.HealthCheckURL, err)
		} else {
			servers.Status = StatusWarmingUp
		}

		servers.Mutex.Unlock()
		
	
	
}

