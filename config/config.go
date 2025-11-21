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
    
    for _, serverConfig := range cfg.Servers {
    	
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
