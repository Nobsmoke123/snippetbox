package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDb(t *testing.T, dsn string) *pgxpool.Pool {

	context := context.Background()
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the "multiStatements=true" parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec() call

	pool, err := pgxpool.New(context, dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from the file and execute the statements, closing
	// the connection pool and calling t.Fatal() in the event of an error.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}

	_, err = pool.Exec(context, string(script))
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}

	// Use t.Cleanup() to register a function *which will automatically be
	// called by Go when the current test (or sub-test) which calls newTestDB()
	// has finished*. In this function we read and execute the teardown script,
	// and close the database connection pool.
	t.Cleanup(func() {
		defer pool.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = pool.Exec(context, string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return pool
}
