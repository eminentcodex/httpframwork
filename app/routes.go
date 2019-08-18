package app

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"httpframwork/app/api"
	"httpframwork/app/middleware"
)

type AppRoutes struct {
	Name    string
	Path    string
	Method  []string
	Handler http.HandlerFunc
}

var (
	allowedOrigin = []string{"*"}

	allowedHeaders = []string{"" +
		"X-Requested-With",
		"Content-Type",
		"X-CustomHeader",
		"Keep-Alive",
		"User-Agent",
		"X-Requested-With",
		"If-Modified-Since",
		"Cache-Control",
		"Authorization",
	}

	allowedMethods = []string{
		"OPTIONS",
		"HEAD",
		"GET",
		"POST",
		"PUT",
		"DELETE",
	}
)

func (r AppRoutes) New(name, path string, methods []string, handler http.HandlerFunc) (*AppRoutes) {

	rt := &AppRoutes{
		Name:    name,
		Path:    path,
		Method:  methods,
		Handler: handler,
	}

	return rt
}

// Prepare routes
func (a *Application) prepareRoutes() http.Handler {

	router := mux.NewRouter()

	// Register the routes here
	a.Routes = []*AppRoutes{
		AppRoutes{}.New(api.RegisterHeartbeat(a.Container, a.Config)),
		//AppRoutes{}.New(api.RegisterSample(a.Container, a.Config)),
	}

	for _, r := range a.Routes {
		router.
			Path(r.Path).
			Name(r.Name).
			HandlerFunc(r.Handler).
			Methods(r.Method...)
	}

	a.initMiddleware(router)
	originsOk := handlers.AllowedOrigins(allowedOrigin)
	headersOk := handlers.AllowedHeaders(allowedHeaders)
	methodsOk := handlers.AllowedMethods(allowedMethods)
	//nrgorilla.InstrumentRoutes(a.Server.Router, a.NewRelic)

	handler := handlers.LoggingHandler(
		a.Log.Writer(),
		handlers.CORS(originsOk, headersOk, methodsOk)(func(m *mux.Router) http.Handler {
			return m
		}(router)), )

	return handler
}

func (a *Application) initMiddleware(router *mux.Router) {
	router.Use(middleware.Sample)
}
