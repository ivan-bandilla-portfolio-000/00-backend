package providers

import (
	"net/http"
)

type RouteProvider struct {
	Mux *http.ServeMux
}

func NewRouteProvider() *RouteProvider {
	return &RouteProvider{Mux: http.NewServeMux()}
}

func (rp *RouteProvider) Register(path string, handler http.HandlerFunc) {
	rp.Mux.HandleFunc(path, handler)
}
