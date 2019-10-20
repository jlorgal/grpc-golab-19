package metrics

import (
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// ServePrometheusEndpoint serves an HTTP endpoint for prometheus metrics of a gRPC service
func ServePrometheusEndpoint(srv *grpc.Server, target string) error {
	grpc_prometheus.Register(srv)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(target, nil)
}
