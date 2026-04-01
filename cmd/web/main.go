package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nobsmoke123/snippetbox/internal/models"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger         *slog.Logger
	users          *models.UserModel
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level:     slog.LevelDebug,
	// 	AddSource: true, // if you wanna add the source to the log
	// }))

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load the data from the .env file
	err := godotenv.Load(".env")

	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	dsn := os.Getenv("DATABASE_URL")

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDb(dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()

	// Initialize the decoder
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our PostgreSQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		users: 			&models.UserModel{DB: db},
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with
	// the * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message.
	// Print a log message to say that the server is starting
	// log.Printf("starting the server on %s", *addr)
	logger.Info("starting server at", slog.String("addr", *addr))

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servermux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	// err = http.ListenAndServe(*addr, app.routes())

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion: tls.VersionTLS13,
	}

	// Initialize a new http.Server struct. We set the addr and the Handler fields
	// so that the server uses the same network address and routes as before\
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the ErrorLog field. If
		// you would prefer to log the server errors at Warn level instead, you
		// could pass slog.LevelWarn as the final parameter.
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// err = srv.ListenAndServe() // listen and serve on HTTP
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// log.Fatal(err)
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())

	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Successfully connected to DB!")
	return db, nil
}
