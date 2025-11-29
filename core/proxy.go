package core

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(w http.ResponseWriter, r *http.Request, target *url.URL){
	// holy shit of a function name
	proxy:= httputil.NewSingleHostReverseProxy(target)

	proxy.ServeHTTP(w,r)
}