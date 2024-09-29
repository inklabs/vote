package postgresrepo

import (
	"fmt"
	"os"
)

// Config holds the state for a PostgreSQL DB config.
type Config struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SearchPath string
}

// DataSourceName returns the DSN for a PostgreSQL DB.
func (c Config) DataSourceName() string {
	searchPath := ""
	if c.SearchPath != "" {
		searchPath = fmt.Sprintf(" search_path=%s", c.SearchPath)
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, searchPath)
}

// NewConfigFromEnvironment loads a Postgres config from environment variables.
func NewConfigFromEnvironment() (Config, error) {
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDBName := os.Getenv("PG_DBNAME")
	pgSearchPath := os.Getenv("PG_SEARCH_PATH")

	if pgHost == "" || pgUser == "" || pgDBName == "" {
		return Config{}, fmt.Errorf("postgreSQL DB has not been configured via environment variables")
	}

	port := "5432"
	if pgPort != "" {
		port = pgPort
	}

	return Config{
		Host:       pgHost,
		Port:       port,
		User:       pgUser,
		Password:   pgPassword,
		DBName:     pgDBName,
		SearchPath: pgSearchPath,
	}, nil
}
