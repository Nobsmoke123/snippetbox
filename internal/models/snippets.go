package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	var id int
	stmt := `INSERT into snippets (title, content, created, expires) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + $3 * INTERVAL '1 day') RETURNING id`

	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	var snippet Snippet
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP AND id=$1`

	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&snippet.ID, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {
		fmt.Println(err.Error())
		return Snippet{}, err
	}
	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
