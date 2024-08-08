package steps

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/service"
	"github.com/ONSdigital/dp-legacy-cache-proxy/service/mock"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type Component struct {
	componenttest.ErrorFeature
	svcList                *service.ExternalServiceList
	svc                    *service.Service
	errorChan              chan error
	Config                 *config.Config
	HTTPServer             *http.Server
	ServiceRunning         bool
	apiFeature             *componenttest.APIFeature
	babbageFeature         *BabbageFeature
	legacyCacheAPIFeature  *LegacyCacheAPIFeature
	releaseCalendarFeature *ReleaseCalendarFeature
}

func NewComponent() (*Component, error) {
	c := &Component{
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	c.babbageFeature = NewBabbageFeature()
	c.legacyCacheAPIFeature = NewLegacyCacheAPIFeature()
	c.releaseCalendarFeature = NewReleaseCalendarFeature()

	c.Config.BabbageURL = c.babbageFeature.Server.URL
	c.Config.LegacyCacheAPIURL = c.legacyCacheAPIFeature.Server.URL
	c.Config.RelCalURL = c.releaseCalendarFeature.Server.URL
	c.Config.EnablePublishExpiryOffset = true

	initMock := &mock.InitialiserMock{
		DoGetHealthCheckFunc:       c.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:        c.DoGetHTTPServer,
		DoGetRequestMiddlewareFunc: c.DoGetRequestMiddleware,
	}

	c.svcList = service.NewServiceList(initMock)

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	c.legacyCacheAPIFeature.Reset()
	return c
}

func (c *Component) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.babbageFeature.Server.Close()
		c.legacyCacheAPIFeature.Server.Close()
		c.releaseCalendarFeature.Server.Close()
		if err := c.svc.Close(context.Background()); err != nil {
			return err
		}
		c.ServiceRunning = false
	}
	return nil
}

var (
	portWobble = 22222
	portMux    sync.RWMutex
)

func getBindAddr() (bindURL string) {
	portMux.Lock()
	portWobble++
	bindURL = "localhost:" + strconv.Itoa(portWobble)
	portMux.Unlock()
	return
}

func (c *Component) InitialiseService() (http.Handler, error) {
	c.Config.BindAddr = getBindAddr()
	var err error
	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, "1", "", "", c.errorChan)
	if err != nil {
		return nil, err
	}

	c.ServiceRunning = true
	return c.HTTPServer.Handler, nil
}

func (c *Component) DoGetHealthcheckOk(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
	// nolint:revive // param names give context here.
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(_ *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer = &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              bindAddr,
		Handler:           router,
	}
	return c.HTTPServer
}

func (c *Component) DoGetRequestMiddleware() service.RequestMiddleware {
	return &HTTPTestRequestMiddleware{}
}

type HTTPTestRequestMiddleware struct{}

func (rm HTTPTestRequestMiddleware) GetMiddlewareFunction() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The APIFeature in dp-component-test appends "http://foo" to the request. In production, the scheme and
			// host are not set. This middleware removes them, so that the request looks like it would in production.
			r.URL.Scheme = ""
			r.URL.Host = ""

			requestURI, _ := strings.CutPrefix(r.RequestURI, "http://foo")
			r.RequestURI = requestURI

			next.ServeHTTP(w, r)
		})
	}
}
