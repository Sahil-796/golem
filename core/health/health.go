

package health

import (
    "time"
    "github.com/Sahil-796/golem/types"
)

func StartHealthCheckers(servers []*types.Server, cfg []types.ServerConfig) {
    for i := range servers {
        interval := cfg[i].HealthCheckConfig.Interval

        go func(i int) {
            ticker := time.NewTicker(interval)
            for range ticker.C {
                ActiveCheckSingle(servers[i], cfg[i]) // single server check
            }
        }(i)  // pass the i and dont use it directly
    } 
}
