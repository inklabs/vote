package postgresrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/internal/electionrepository/postgresrepo"
)

func TestConfig(t *testing.T) {
	t.Run("NewConfigFromEnvironment", func(t *testing.T) {
		t.Run("errors when not set", func(t *testing.T) {
			// Given
			key := "PG_HOST"
			t.Setenv(key, "")

			// When
			_, err := postgresrepo.NewConfigFromEnvironment()

			// Then
			require.EqualError(t, err, "postgreSQL DB has not been configured via environment variables")
		})

		t.Run("returns correct DSN", func(t *testing.T) {
			// Given
			t.Setenv("PG_HOST", "host")
			t.Setenv("PG_PORT", "8080")
			t.Setenv("PG_USER", "user")
			t.Setenv("PG_PASSWORD", "password")
			t.Setenv("PG_DBNAME", "dbname")
			t.Setenv("PG_SEARCH_PATH", "searchpath")

			// When
			config, err := postgresrepo.NewConfigFromEnvironment()
			require.NoError(t, err)

			// Then
			dsn := config.DataSourceName()
			assert.Equal(t, "host=host port=8080 user=user password=password dbname=dbname sslmode=disable search_path=searchpath", dsn)
		})
	})

	t.Run("DataSourceName", func(t *testing.T) {
		t.Run("without search path", func(t *testing.T) {
			// Given
			config := &postgresrepo.Config{
				Host:     "host",
				Port:     "8080",
				User:     "user",
				Password: "password",
				DBName:   "dbname",
			}

			// When
			dsn := config.DataSourceName()

			// Then
			assert.Equal(t, "host=host port=8080 user=user password=password dbname=dbname sslmode=disable", dsn)
		})

		t.Run("all values", func(t *testing.T) {
			// Given
			config := &postgresrepo.Config{
				Host:       "host",
				Port:       "8080",
				User:       "user",
				Password:   "password",
				DBName:     "dbname",
				SearchPath: "searchpath",
			}

			// When
			dsn := config.DataSourceName()

			// Then
			assert.Equal(t, "host=host port=8080 user=user password=password dbname=dbname sslmode=disable search_path=searchpath", dsn)
		})
	})
}
