package postgre

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"

	"supportflow/core"
)

var Pool *pgxpool.Pool

func Init(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%d",
		core.GetString("db.postgres.username", "postgres"),
		core.GetString("db.postgres.password", "postgres"),
		core.GetString("db.postgres.host", "localhost"),
		core.GetString("db.postgres.port", "5432"),
		core.GetString("db.postgres.database", "supportflow"),
		core.GetString("db.postgres.sslmode", "disable"),
		core.GetInt("db.postgres.maxcon", 10),
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("postgre: connect: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("postgre: ping: %w", err)
	}

	Pool = pool
	return nil
}

func RunMigrations(ctx context.Context) error {
	dir := core.GetString("db.migrations.path", "db/postgre/migrations")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("postgre: read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return fmt.Errorf("postgre: read %s: %w", f, err)
		}
		if _, err := Pool.Exec(ctx, string(data)); err != nil {
			return fmt.Errorf("postgre: exec %s: %w", f, err)
		}
		fmt.Printf("Migration applied: %s\n", f)
	}
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
