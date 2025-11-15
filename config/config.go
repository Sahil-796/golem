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
    
    list, ok := rawTargets.([]any)
    if !ok {
    	return nil, fmt.Errorf("targets must be a list")
    }
    
    for _, item := range list {
    
    	m, ok := item.(map[string]any)
     
     	if !ok {
      		return nil, fmt.Errorf("target urls must have valid syntax")
      	}
     	rawURL := m["url"].(string)
      
       parsed, err := url.Parse(rawURL)
       if err != nil {
          return nil, fmt.Errorf("invalid url '%s': %w", rawURL, err)
       }
      	
       cfg.Servers = append(cfg.Servers, types.Server{
       		URL: parsed, 
         	IsHealthy: true})
       
    }
    
    return cfg, nil
}
