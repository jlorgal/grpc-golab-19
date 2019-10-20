package tracer

import (
	opentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

// NewZipkinTracer creates a zipkin tracer with type opentracing.Tracer
func NewZipkinTracer(reporterURL string, serviceName string, serviceTarget string) (opentracing.Tracer, error) {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(reporterURL)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, serviceTarget)
	if err != nil {
		return nil, err
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.InitGlobalTracer(tracer)

	return tracer, nil
}
