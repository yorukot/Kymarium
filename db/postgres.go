package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Register postgres migrations driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Register file source for migrations.
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/utils/id"
	"go.uber.org/zap"
)

// InitDatabase initialize the database connection pool and return the pool and also migrate the database
func InitDatabase() (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Configure connection pool to handle concurrent operations better
	config, err := pgxpool.ParseConfig(getDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Increase pool size to handle more concurrent connections
	config.MaxConns = 25
	config.MinConns = 5

	// Reduce prepared statement cache to prevent "conn busy" errors
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	zap.L().Info("Database initialized")

	runAll := len(os.Args) < 2 || os.Args[1] == "all"
	if runAll || os.Args[1] == "api" {
		Migrator()
		if err := createRegionsDataIfNotExists(ctx, pool); err != nil {
			pool.Close()
			return nil, err
		}
	}

	
	
	return pool, nil
}

// getDatabaseURL return a pgsql connection uri by the environment variables
func getDatabaseURL() string {
	dbHost := config.Env().DBHost
	dbPort := config.Env().DBPort
	dbUser := config.Env().DBUser
	dbPassword := config.Env().DBPassword
	dbName := config.Env().DBName
	dbSSLMode := config.Env().DBSSLMode
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)
}

// Migrator the database
func Migrator() {
	zap.L().Info("Migrating database")

	wd, _ := os.Getwd()

	databaseURL := getDatabaseURL()
	migrationsPath := "file://" + wd + "/migrations"

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		zap.L().Fatal("failed to create migrator", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.L().Fatal("failed to migrate database", zap.Error(err))
	}

	zap.L().Info("Database migrated")
}

func createRegionsDataIfNotExists(ctx context.Context, pool *pgxpool.Pool) error {
	zap.L().Info("Creating regions data if not exists")

	iso3166_2 := regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)

	for _, raw := range config.Env().AppRegions {
		name := strings.ToUpper(strings.TrimSpace(raw))
		if name == "" {
			continue
		}
		if !iso3166_2.MatchString(name) {
			return fmt.Errorf("invalid region format %q (expected ISO 3166-2 like US-CA, TW-TPE)", raw)
		}

		regionID, err := id.GetID()
		if err != nil {
			return fmt.Errorf("generate region id: %w", err)
		}

		if _, err := pool.Exec(ctx, `
			INSERT INTO regions (id, name)
			VALUES ($1, $2)
			ON CONFLICT (name) DO NOTHING
		`, regionID, name); err != nil {
			return fmt.Errorf("insert region %q: %w", name, err)
		}
	}

	zap.L().Info("Regions data created or already exists")
	return nil
}
