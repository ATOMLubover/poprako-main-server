package http

import (
	"strconv"

	"poprako-main-server/internal/state"

	"github.com/kataras/iris/v12"
)

func Run(appState state.AppState) {
	// A default Iris application,
	// with logger and recovery middleware already attached.
	app := iris.Default()

	// `appState` is injected here as a dependency.
	// Every handler below can receive it as a parameter.
	//
	// NOTE: the type of `appState` must be a pointer type,
	// as it will then be shared across all running handlers.
	app.RegisterDependency(&appState)

	// Apply handlers to routes.
	routeApp(app)

	// Run the application.
	runServer(app, appState.Cfg.Port)
}

// Route application endpoints.
func routeApp(app *iris.Application) {
}

func runServer(app *iris.Application, port uint16) {
	// TODO: make app runs at 0.0.0.0.
	addr := ":" + strconv.Itoa(int(port))

	app.Listen(addr)
}
