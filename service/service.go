package service

import (
	"context"
	"net"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/proxy"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/netutil"
)

// Service contains all the configs, server and clients to run the proxy
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	Proxy       *proxy.Proxy
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	router := mux.NewRouter()

	var server HTTPServer

	if cfg.OtelEnabled {
		otelHandler := otelhttp.NewHandler(router, "/")
		router.Use(otelmux.Middleware(cfg.OTServiceName))
		server = serviceList.GetHTTPServer(cfg, cfg.BindAddr, otelHandler)
	} else {
		server = serviceList.GetHTTPServer(cfg, cfg.BindAddr, router)
	}

	// TODO: Any middleware will require 'otelhttp.NewMiddleware(cfg.OTServiceName),' included for Open Telemetry
	router.Use(serviceList.Init.DoGetRequestMiddleware().GetMiddlewareFunction())

	// TODO: Add other(s) to serviceList here

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)

	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	router.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	// The proxy needs to be set up after the HealthCheck route has been added to the router: in the Setup method, the
	// proxy adds a catch-all route, so any other routes added after that one will never be reachable.
	p := proxy.Setup(ctx, router, cfg)

	hc.Start(ctx)

	// Create a LimitListener to cap concurrent http connections
	l, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		log.Fatal(ctx, "error starting tcp listener", err)
	}

	if maxC := cfg.HTTPMaxConnections; maxC > 0 {
		l = netutil.LimitListener(l, maxC)
	}

	// Run the http server in a new go-routine
	go func(l net.Listener) {
		if err := server.Serve(l); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
		defer l.Close()
	}(l)

	return &Service{
		Config:      cfg,
		Router:      router,
		Proxy:       p,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      server,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(_ context.Context, _ HealthChecker) (err error) {
	// TODO: add other health checks here, as per dp-upload-service

	return nil
}
