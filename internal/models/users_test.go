package models

import (
	"os"
	"testing"

	"github.com/Nobsmoke123/snippetbox/internal/assert"
	"github.com/joho/godotenv"
)

func TestUserModelExists(t *testing.T) {
	file, err := os.Open("./../../.env")
	if err != nil {
		file.Close()
		t.Fatal(err)
	}

	defer file.Close()

	err = godotenv.Load(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	test_database := os.Getenv("TEST_DATABASE_URL")

	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 2,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test.
			db := newTestDb(t, test_database)

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test.
			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
