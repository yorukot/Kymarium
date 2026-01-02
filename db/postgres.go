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
	"github.com/yorukot/kymarium/models"
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
		if err := creteRegionsDataIfNotExists(ctx, pool); err != nil {
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

func creteRegionsDataIfNotExists(ctx context.Context, pool *pgxpool.Pool) error {
	
	zap.L().Info("Creating regions data if not exists")

	// Check if regions already exist
	regions := make([]models.Region, 0, len(config.Env().AppRegions))
	for _, raw := range config.Env().AppRegions {
		parts := strings.SplitN(raw, "-", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid region format %q (expected CC-City)", raw)
		}

		countryCode := strings.TrimSpace(parts[0])
		city := prettifyCity(strings.TrimSpace(parts[1]))

		countryName, ok := countryNames[countryCode]
		if !ok {
			countryName = countryCode
		}

		displayName := fmt.Sprintf("%s, %s", countryName, city)
		regionID, err := id.GetID()
		if err != nil {
			return fmt.Errorf("generate region id: %w", err)
		}

		regions = append(regions, models.Region{
			ID:          regionID,
			Name:        raw,
			DisplayName: displayName,
		})
	}

	for _, region := range regions {
		if _, err := pool.Exec(ctx, `
			INSERT INTO regions (id, name, display_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (name) DO NOTHING
		`, region.ID, region.Name, region.DisplayName); err != nil {
			return fmt.Errorf("insert region %q: %w", region.Name, err)
		}
	}
	
	zap.L().Info("Regions data created or already exists")
	return nil
}

var camelCaseWordSplit = regexp.MustCompile(`([a-z])([A-Z])`)

func prettifyCity(city string) string {
	city = camelCaseWordSplit.ReplaceAllString(city, `$1 $2`)
	return strings.TrimSpace(city)
}

var countryNames = map[string]string{
	"TW": "Taiwan",
	"US": "United States",
	"UK": "United Kingdom",
	"CA": "Canada",
	"SG": "Singapore",
	"JP": "Japan",
	"KR": "South Korea",
	"AU": "Australia",
	"IN": "India",
	"DE": "Germany",
	"FR": "France",
}
