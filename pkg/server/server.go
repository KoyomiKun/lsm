package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	uuid "github.com/satori/go.uuid"

	"lsm/pkg/logger"
)

var (
	mtxDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "lsm_db_http_duration",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path", "method", "code"})

	mtxReqSize = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "lsm_db_http_request_bytes",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path", "method", "code"})

	mtxResSize = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "lsm_db_http_response_bytes",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path", "method", "code"})
)

type Server interface {
	Serve()
}

type HTTPServer struct {
	*http.Server

	mux    *http.ServeMux
	logger logger.Logger
}

func NewHTTPServer(addr string) *HTTPServer {
	mux := http.NewServeMux()
	return &HTTPServer{
		&http.Server{
			Addr:    addr,
			Handler: mux,
		},
		mux,
		logger.GetGlobalLogger().WithModule("server"),
	}
}

func (hs *HTTPServer) Serve(stop <-chan struct{}) {
	go func() {
		if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hs.logger.Fatal().AnErr("err", err).Msg("listen and serve failed")
		}
	}()
	<-stop
	if err := hs.Shutdown(context.Background()); err != nil {
		hs.logger.Fatal().AnErr("err", err).Msg("shutdown server failed")
	}

}

func (hs *HTTPServer) Register(ctx context.Context, path string, handler http.Handler) {
	wrapCtx := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.NewV4().String()
			handler.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "RequestID", reqID)))
		})
	}

	handler = wrapCtx(handler)

	cm := map[string]string{"path": path}
	handler = promhttp.InstrumentHandlerDuration(mtxDuration.MustCurryWith(cm), handler)
	handler = promhttp.InstrumentHandlerRequestSize(mtxReqSize.MustCurryWith(cm), handler)
	handler = promhttp.InstrumentHandlerResponseSize(mtxResSize.MustCurryWith(cm), handler)
	hs.mux.Handle(path, handler)
}
