package main

import (
	"elastic/handler"
	"elastic/l"
	"elastic/store"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	render_chi "github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// func initJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
func initJaeger(service string, logger l.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(logger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {

	logger, err := l.CreateZapLogger()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	tracer, closer := initJaeger("example", logger)
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	//Initialize Stores
	articleStore, err := store.NewArticleStore()
	parseErr(err, &logger)
	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore, tracer, logger)

	// chi
	r := chi.NewRouter()
	r.Use(render_chi.SetContentType(render_chi.ContentTypeJSON))

	//routes
	r.Get("/article/id/{id}", articleHandler.Id_chi)
	r.Post("/article/add", articleHandler.Add_chi)
	r.Post("/article/search", articleHandler.Search_chi)
	//panic
	panicHandler := handler.PanicHandler{}
	r.Get("/panic", panicHandler.Handle_chi)
	r.Post("/log/add", panicHandler.Log_chi)

	http.ListenAndServe(":3333", r)
}

func parseErr(err error, logger *l.Logger) {
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Infof("Application started")
}
