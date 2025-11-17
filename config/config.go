package config

import (
	"fmt"
	"net/url"

	"github.com/Sahil-796/golem/types"
	"github.com/spf13/viper"
)

func LoadConfig() (*types.Config, error){
	
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err!=nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	
	cfg := &types.Config{}
	
	cfg.Strategy = viper.GetString("loadbalancer.strategy")
	rawTargets := viper.Get("targets")
	
	if rawTargets == nil {
		return nil, fmt.Errorf("no targets provided")
    }
    
    // asserting that rawTargets is an empty list with any types in it
    list, ok := rawTargets.([]any) 
    if !ok {
    	return nil, fmt.Errorf("targets must be a list")
    }
    
    for _, item := range list {
    	// asserting again
     	// the yaml returns url: ""  
    	m, ok := item.(map[string]any)
     	// if !ok, the cinfig,yaml syntax might have issues.  
     	if !ok {
      		return nil, fmt.Errorf("target urls must have valid syntax")
      	}
     	rawURL := m["url"].(string)
      
      //parsing url from config.yaml to net/url
       parsed, err := url.Parse(rawURL)
       if err != nil {
          return nil, fmt.Errorf("invalid url '%s': %w", rawURL, err)
       }
      	
       cfg.Servers = append(cfg.Servers, types.Server{
       		URL: parsed, 
        	Threshold: m["threshold"].(int),
         	healthEndpoint: m["healthEndpoint"]
         	IsHealthy: true,
       }) //mutex - auto initialised
       
    }
    
    return cfg, nil
}
