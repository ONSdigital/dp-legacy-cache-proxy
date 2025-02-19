package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/service"
	"github.com/ONSdigital/dp-legacy-cache-proxy/service/mock"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"

	errServer      = errors.New("HTTP Server error")
	errHealthcheck = errors.New("healthCheck error")

	bindAddrAny = "localhost:0"
)

// nolint:revive // param names give context here.
var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

// nolint:revive // param names give context here.
var funcDoGetHTTPServerNil = func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

func TestRun(t *testing.T) {
	Convey("Having a correctly initialised service and set of mocked dependencies", t, func() {
		cfg, cfgErr := config.Get()
		So(cfgErr, ShouldBeNil)

		cfg.BindAddr = bindAddrAny

		// nolint:revive // param names give context here.
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		failingServerMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		// nolint:revive // param names give context here.
		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		// nolint:revive // param names give context here.
		funcDoGetHTTPServer := func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		// nolint:revive // param names give context here.
		funcDoGetFailingHTTPServer := func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
			return failingServerMock
		}

		funcDoGetRequestMiddleware := func() service.RequestMiddleware {
			return &service.NoOpRequestMiddleware{}
		}

		Convey("Given that initialising healthcheck returns an error", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc:        funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc:       funcDoGetHealthcheckErr,
				DoGetRequestMiddlewareFunc: funcDoGetRequestMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.HealthCheck, ShouldBeFalse)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc:        funcDoGetHTTPServer,
				DoGetHealthCheckFunc:       funcDoGetHealthcheckOk,
				DoGetRequestMiddlewareFunc: funcDoGetRequestMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.HealthCheck, ShouldBeTrue)
			})

			Convey("The checkers are registered and the healthcheck and http server started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 0)
				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Convey("Then the proxy's catch-all route is the one with the lowest precedence", func() {
				var lastRouteName string
				// nolint:revive // param names give context here.
				_ = svc.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
					lastRouteName = route.GetName()
					return nil
				})
				So(lastRouteName, ShouldEqual, "Proxy Catch-All")
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		/* ADD CODE OR REMOVE: put this code in, if you have Checkers to register
		Convey("Given that Checkers cannot be registered", func() {
			// setup (run before each `Convey` at this scope / indentation):
			errAddheckFail := errors.New("Error(s) registering checkers for healthcheck")
			hcMockAddFail := &mock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddheckFail },
				StartFunc:    func(ctx context.Context) {},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				// ADD CODE: add the checkers that you want to register here
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				// ADD CODE: add code to confirm checkers exist
				So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 0) // ADD CODE: change the '0' to the number of checkers you have registered
			})
			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})*/

		Convey("Given that all dependencies are successfully initialised but the http server fails", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHealthCheckFunc:       funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:        funcDoGetFailingHTTPServer,
				DoGetRequestMiddlewareFunc: funcDoGetRequestMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			Convey("Then the error is returned in the error channel", func() {
				sErr := <-svcErrors
				So(sErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		cfg, cfgErr := config.Get()
		So(cfgErr, ShouldBeNil)

		cfg.BindAddr = bindAddrAny
		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		// nolint:revive // param names give context here.
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		// nolint:revive // param names give context here.
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			// nolint:revive // param names give context here.
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetRequestMiddlewareFunc: func() service.RequestMiddleware { return &service.NoOpRequestMiddleware{} },
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			// nolint:revive // param names give context here.
			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}
			// nolint:revive // param names give context here.
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
					return failingserverMock
				},
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetRequestMiddlewareFunc: func() service.RequestMiddleware { return &service.NoOpRequestMiddleware{} },
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "failed to shutdown gracefully")
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			// nolint:revive // param names give context here.
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(200 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err := svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "context deadline exceeded")
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
