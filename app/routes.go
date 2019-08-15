package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"httpframwork/app/api"
	"httpframwork/app/middleware"
)

type AppRoutes struct {
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func (r AppRoutes) New(name, path, method string, handler http.HandlerFunc) (*AppRoutes) {

	rt := &AppRoutes{
		Name:    name,
		Path:    path,
		Method:  method,
		Handler: handler,
	}

	return rt
}

// Prepare routes
func (a *Application) prepareRoutes() *mux.Router {

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
			Methods(r.Method)
	}

	a.initMiddleware(router)

	//nrgorilla.InstrumentRoutes(a.Server.Router, a.NewRelic)

	return router
}

func (a *Application) initMiddleware(router *mux.Router) {
	router.Use(middleware.Sample)
}
