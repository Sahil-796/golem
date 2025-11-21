package core

import (
	"net/http"
	"net/url"

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
	
	for i := range servers {
		
		baseUrl := servers[i].URL
		health := baseUrl.ResolveReference(&url.URL{Path: ServerConfig[i].HealthCheckConfig.Path})
		result, err := http.Get(health.String())
		
		if result != nil {
			defer result.Body.Close()
		}
		
		servers[i].Mutex.Lock()
		if err != nil || result.StatusCode != http.StatusOK {
			
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

