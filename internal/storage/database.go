package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ MetricsStorage = DatabaseStorage{}

// implements pgxpool.Pool
type PGXPool interface {
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

type DatabaseStorage struct {
	pool PGXPool
}

func NewDatabaseStorage(dsn string) (*DatabaseStorage, error) {
	migrator := NewMigrator(dsn, "file://db/migrate", 5)

	if err := migrator.Run(); err != nil {
		return nil, fmt.Errorf("migrations run failed: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool init failed: %w", err)
	}

	return &DatabaseStorage{pool: pool}, nil
}

func (d DatabaseStorage) Push(ctx context.Context, key string, record Record) error {
	sql := "INSERT INTO metrics(id, name, kind, value) values ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET value = $4"
	_, err := d.pool.Exec(ctx, sql, key, record.Name, record.Value.Kind(), record.Value.String())

	if err != nil {
		return fmt.Errorf("db storage Push() error: %w", err)
	}

	return nil
}

func (d DatabaseStorage) Get(ctx context.Context, key string) (Record, error) {
	var (
		name   string
		kind   string
		value  float64
		record Record
		err    error
	)

	sql := "SELECT name, kind, value FROM metrics WHERE id=$1"
	err = d.pool.QueryRow(ctx, sql, key).Scan(&name, &kind, &value)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Record{}, entities.ErrRecordNotFound
		}

		return record, fmt.Errorf("db storage Get() error: %w", err)
	}

	switch kind {
	case metrics.KindCounter:
		record, err = Record{Name: name, Value: metrics.Counter(value)}, nil
	case metrics.KindGauge:
		record, err = Record{Name: name, Value: metrics.Gauge(value)}, nil
	default:
		err = fmt.Errorf("db storage kind=%s unknown", kind)
	}

	return record, err
}

func (d DatabaseStorage) List(ctx context.Context) ([]Record, error) {
	rows, err := d.pool.Query(ctx, "SELECT name, kind, value FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("db storage List() error: %w", err)
	}

	defer rows.Close()

	var (
		name  string
		kind  string
		value float64
	)

	result := make([]Record, 0)
	_, err = pgx.ForEachRow(rows, []any{&name, &kind, &value}, func() error {
		switch kind {
		case metrics.KindCounter:
			result = append(result, Record{Name: name, Value: metrics.Counter(value)})
			return nil

		case metrics.KindGauge:
			result = append(result, Record{Name: name, Value: metrics.Gauge(value)})
			return nil

		default:
			return fmt.Errorf("db storage kind=%s unknown", kind)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("db storage List() error: %w", err)
	}

	return result, nil
}

func (d DatabaseStorage) Ping(ctx context.Context) error {
	if err := d.pool.Ping(ctx); err != nil {
		return fmt.Errorf("db storage Ping() error: %w", err)
	}

	return nil
}

func (d DatabaseStorage) Close(ctx context.Context) error {
	d.pool.Close()
	return nil
}
