package main

import (
	"net/http"

	"github.com/Nobsmoke123/snippetbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new ServeMux,
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	// fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Add new GET /ping route.
	mux.HandleFunc("GET /ping", ping)

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Register the other application routes as normal
	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))

	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignUp))
	
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignUpPost))

	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	mux.Handle("GET /account/settings", protected.ThenFunc(app.settingsPage))

	mux.Handle("POST /account/settings", protected.ThenFunc(app.settingsPagePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
