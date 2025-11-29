package strategy

import (
	"github.com/Sahil-796/golem/types"
)

// defining an interface to implement the same method for types having the same method
type Strategy interface {
	Next(servers []*types.Server) *types.Server
}

func Get(name string) Strategy {
    switch name {
    case "round_robin":
        return &RoundRobin{}
    default:
        return &RoundRobin{}
    }
}