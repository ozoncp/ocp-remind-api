package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/ozoncp/ocp-remind-api/internal/configuration"
)

// InitTracer - init tracer
func InitTracer(serviceName string) io.Closer {
	cfgMetrics := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: configuration.Instance().Jaeger.Host + ":" +
				configuration.Instance().Jaeger.Port,
		},
	}
	tracer, closer, err := cfgMetrics.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Err(err).Msgf("failed init jaeger: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
	log.Info().Msgf("Traces started")

	return closer
}
