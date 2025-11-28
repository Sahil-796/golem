package strategy

import (
	"github.com/Sahil-796/golem/types"
)

// defining an interface to implement the same method for types having the same method
type Strategy interface {
	Next(servers []*types.Server) *types.Server
}

