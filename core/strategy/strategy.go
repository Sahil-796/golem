package strategy

import (
	"github.com/Sahil-796/golem/types"
	"net/http"
)

// defining an interface to implement the same method for types having the same method
type Strategy interface {
	Next(r *http.Request, servers []*types.Server) *types.Server
}

func Get(name string) Strategy {
	switch name {
	case "round_robin":
		return &RoundRobin{index: -1}
	case "weighted_round_robin":
		return &WeightedRoundRobin{index: 0}
	case "least_connections":
		return &LeastConnections{}
	case "ip_hash":
		return &IPHash{}
	default:
		return &RoundRobin{}
	}
}
