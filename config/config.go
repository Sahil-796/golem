package config

import (
	"fmt"
	"net/url"

	"github.com/Sahil-796/golem/types"
	"github.com/spf13/viper"
)


// change this according to removed url and keep a final - parsed, built, validated url 
// in Server.URL the concurrent one.

func LoadConfig() (*types.Config, []*types.Server, error){
	
	// setup of viper
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err!=nil {
		return nil, nil, fmt.Errorf("error reading config: %w", err)
	}
	
	
	cfg := &types.Config{} 
	
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
    	
        // parsing a final url and saving it in Server.URL url.URL
    	healthURL, err := BuildFinalURL(serverConfig.HealthCheckConfig)
     	if err!=nil {
      		return nil, nil, fmt.Errorf("invalid url : %w", err)
      	}
       
       baseURL := &url.URL{
           Scheme: healthURL.Scheme,
           Host:   healthURL.Host,
       }
       
       	server := &types.Server {
	       URL: baseURL,
		   HealthCheckURL: healthURL,							
	       IsHealthy: true,
		   Status: "initial",
	       ConsecutiveFailures: 0,
	       ConsecutiveSuccesses: 0,
        }
        
        runtimeServers = append(runtimeServers, server)

    }
    
    return cfg, runtimeServers, nil
}


func BuildFinalURL(hc types.HealthCheckConfig) (*url.URL, error) {
    if hc.Host == "" {
        return nil, fmt.Errorf("cannot build URL: host missing")
    }

    if hc.Protocol != "http" && hc.Protocol != "https" {
        return nil, fmt.Errorf("build url: protocol '%s' not supported yet", hc.Protocol)
    }
    
    host := hc.Host
       if hc.Port > 0 {
           host = fmt.Sprintf("%s:%d", hc.Host, hc.Port)
       }
   
    built := &url.URL{
           Scheme: hc.Protocol, // http / https
           Host:   host,        // "api:3000" OR "localhost:8080"
           Path:   "/",     // "/health" OR "/"
       }

    return built, nil
}