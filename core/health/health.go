package core

import (
	"github.com/Sahil-796/golem/types"
	"net/http"
	"net/url"
)

func ActiveCheck(servers []*types.Server) {
	
	for i := range servers {
		
		u, _ := servers[i].URL.Parse("http://localhost:3000/")
		health := u.ResolveReference(&url.URL{Path: "/health"})
		result, err := http.Get(health.String())
		
		if err != nil {
			
		}
		
		
		
	}
}

