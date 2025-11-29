package core

import (
	"sync"
	"github.com/Sahil-796/golem/core/strategy"
	"github.com/Sahil-796/golem/types"
)


type LoadBalancer struct {
	Mutex   sync.Mutex 
	Backends []*types.Server
	Strategy strategy.Strategy
}


func NewLoadBalancer (strategyName string, backends []*types.Server) *LoadBalancer {
	return &LoadBalancer{
		Strategy: strategy.Get(strategyName),
		Backends: backends,
	}
}
