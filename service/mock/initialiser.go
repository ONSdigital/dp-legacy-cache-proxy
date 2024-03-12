// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/service"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
//	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
//		// make and configure a mocked service.Initialiser
//		mockedInitialiser := &InitialiserMock{
//			DoGetHTTPServerFunc: func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
//				panic("mock out the DoGetHTTPServer method")
//			},
//			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
//				panic("mock out the DoGetHealthCheck method")
//			},
//			DoGetRequestMiddlewareFunc: func() service.RequestMiddleware {
//				panic("mock out the DoGetRequestMiddleware method")
//			},
//		}
//
//		// use mockedInitialiser in code that requires service.Initialiser
//		// and then make assertions.
//
//	}
type InitialiserMock struct {
	// DoGetHTTPServerFunc mocks the DoGetHTTPServer method.
	DoGetHTTPServerFunc func(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer

	// DoGetHealthCheckFunc mocks the DoGetHealthCheck method.
	DoGetHealthCheckFunc func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error)

	// DoGetRequestMiddlewareFunc mocks the DoGetRequestMiddleware method.
	DoGetRequestMiddlewareFunc func() service.RequestMiddleware

	// calls tracks calls to the methods.
	calls struct {
		// DoGetHTTPServer holds details about calls to the DoGetHTTPServer method.
		DoGetHTTPServer []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// DoGetHealthCheck holds details about calls to the DoGetHealthCheck method.
		DoGetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BuildTime is the buildTime argument value.
			BuildTime string
			// GitCommit is the gitCommit argument value.
			GitCommit string
			// Version is the version argument value.
			Version string
		}
		// DoGetRequestMiddleware holds details about calls to the DoGetRequestMiddleware method.
		DoGetRequestMiddleware []struct {
		}
	}
	lockDoGetHTTPServer        sync.RWMutex
	lockDoGetHealthCheck       sync.RWMutex
	lockDoGetRequestMiddleware sync.RWMutex
}

// DoGetHTTPServer calls DoGetHTTPServerFunc.
func (mock *InitialiserMock) DoGetHTTPServer(cfg *config.Config, bindAddr string, router http.Handler) service.HTTPServer {
	if mock.DoGetHTTPServerFunc == nil {
		panic("InitialiserMock.DoGetHTTPServerFunc: method is nil but Initialiser.DoGetHTTPServer was just called")
	}
	callInfo := struct {
		Cfg      *config.Config
		BindAddr string
		Router   http.Handler
	}{
		Cfg:      cfg,
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockDoGetHTTPServer.Lock()
	mock.calls.DoGetHTTPServer = append(mock.calls.DoGetHTTPServer, callInfo)
	mock.lockDoGetHTTPServer.Unlock()
	return mock.DoGetHTTPServerFunc(cfg, bindAddr, router)
}

// DoGetHTTPServerCalls gets all the calls that were made to DoGetHTTPServer.
// Check the length with:
//
//	len(mockedInitialiser.DoGetHTTPServerCalls())
func (mock *InitialiserMock) DoGetHTTPServerCalls() []struct {
	Cfg      *config.Config
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		Cfg      *config.Config
		BindAddr string
		Router   http.Handler
	}
	mock.lockDoGetHTTPServer.RLock()
	calls = mock.calls.DoGetHTTPServer
	mock.lockDoGetHTTPServer.RUnlock()
	return calls
}

// DoGetHealthCheck calls DoGetHealthCheckFunc.
func (mock *InitialiserMock) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	if mock.DoGetHealthCheckFunc == nil {
		panic("InitialiserMock.DoGetHealthCheckFunc: method is nil but Initialiser.DoGetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}{
		Cfg:       cfg,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		Version:   version,
	}
	mock.lockDoGetHealthCheck.Lock()
	mock.calls.DoGetHealthCheck = append(mock.calls.DoGetHealthCheck, callInfo)
	mock.lockDoGetHealthCheck.Unlock()
	return mock.DoGetHealthCheckFunc(cfg, buildTime, gitCommit, version)
}

// DoGetHealthCheckCalls gets all the calls that were made to DoGetHealthCheck.
// Check the length with:
//
//	len(mockedInitialiser.DoGetHealthCheckCalls())
func (mock *InitialiserMock) DoGetHealthCheckCalls() []struct {
	Cfg       *config.Config
	BuildTime string
	GitCommit string
	Version   string
} {
	var calls []struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}
	mock.lockDoGetHealthCheck.RLock()
	calls = mock.calls.DoGetHealthCheck
	mock.lockDoGetHealthCheck.RUnlock()
	return calls
}

// DoGetRequestMiddleware calls DoGetRequestMiddlewareFunc.
func (mock *InitialiserMock) DoGetRequestMiddleware() service.RequestMiddleware {
	if mock.DoGetRequestMiddlewareFunc == nil {
		panic("InitialiserMock.DoGetRequestMiddlewareFunc: method is nil but Initialiser.DoGetRequestMiddleware was just called")
	}
	callInfo := struct {
	}{}
	mock.lockDoGetRequestMiddleware.Lock()
	mock.calls.DoGetRequestMiddleware = append(mock.calls.DoGetRequestMiddleware, callInfo)
	mock.lockDoGetRequestMiddleware.Unlock()
	return mock.DoGetRequestMiddlewareFunc()
}

// DoGetRequestMiddlewareCalls gets all the calls that were made to DoGetRequestMiddleware.
// Check the length with:
//
//	len(mockedInitialiser.DoGetRequestMiddlewareCalls())
func (mock *InitialiserMock) DoGetRequestMiddlewareCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockDoGetRequestMiddleware.RLock()
	calls = mock.calls.DoGetRequestMiddleware
	mock.lockDoGetRequestMiddleware.RUnlock()
	return calls
}
