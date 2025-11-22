package config

import (
	"fmt"
	"net/url"

	"github.com/Sahil-796/golem/types"
	"github.com/spf13/viper"
)

func LoadConfig() (*types.Config, []*types.Server, error){
	
	// setup of viper
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err!=nil {
		return nil, nil, fmt.Errorf("error reading config: %w", err)
	}
	
	
	cfg := &types.Config{} //return this
	
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, nil, fmt.Errorf("error reading config.yaml: %w", err)
	}

	runtimeServers := make([]*types.Server, 0, len(cfg.Servers))
    
    for i := range cfg.Servers {
    
    	serverConfig := &cfg.Servers[i] //pointer to the server config
    	
    	serverConfig.HealthCheckConfig.SetDefaults() // setting default values
     
     	// validating 
        if err := serverConfig.Validate(); err != nil {
            return nil, nil, fmt.Errorf("invalid server config: %w", err)
        }
    	
    	parsedURL, err := url.Parse(serverConfig.URL)
     	if err!=nil || parsedURL.Host == "" {
      		return nil, nil, fmt.Errorf("invalid url '%s': %w", serverConfig.URL, err)
      	}
       
       	server := &types.Server {
	       URL: parsedURL,
	       IsHealthy: true,
		   Status: "initial",
	       ConsecutiveFailures: 0,
	       ConsecutiveSuccesses: 0,
        }
        
        runtimeServers = append(runtimeServers, server)

    }
    
    return cfg, runtimeServers, nil
}

func buildHealthURL(sc *types.ServerConfig) (*url.URL, error) {
	
	parsed, err := url.Parse(sc.URL)
	if err != nil {
		return nil, fmt.Errorf("Invalid URL %s: %w", sc.URL, err)
	}
	
	// behold generational if statements
	
	if sc.HealthCheckConfig.Protocol != "" {
		parsed.Scheme = sc.HealthCheckConfig.Protocol
	} else {
		parsed.Scheme = "http"
	}
	
	if sc.HealthCheckConfig.Port != "" {
		parsed.Host = fmt.Sprintf("%s:%d")
	} else {
		parsed.Scheme = "http"
	}
	
	return parsed, nil
	
}
