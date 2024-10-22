package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/pkg/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ MetricsStorage = DatabaseStorage{}

// DB Storage
type DatabaseStorage struct {
	Pool PGXPool
}

type dbQueryTracer struct {
	logger *zerolog.Logger
}

func (tracer *dbQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.logger.Debug().Msg(fmt.Sprintf("Executing command \"%s\" with args %v", data.SQL, data.Args))
	return ctx
}

func (tracer *dbQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// empty
}

// DatabseStorage constructor
func NewDatabaseStorage(dsn string) (*DatabaseStorage, error) {
	migrator := NewDatabaseMigrator(dsn, "file://db/migrate", 5)

	if err := migrator.Run(); err != nil {
		return nil, fmt.Errorf("migrations run failed: %w", err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool parsse config failed: %w", err)
	}

	ctx := context.Background()
	config.ConnConfig.Tracer = &dbQueryTracer{logger: &log.Logger}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool init failed: %w", err)
	}

	return &DatabaseStorage{Pool: pool}, nil
}

// Push record to storage
func (d DatabaseStorage) Push(ctx context.Context, key string, record Record) error {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db storage Push() -> Begin() error: %w", err)
	}

	sql := "INSERT INTO metrics(id, name, kind, value) values ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET value = $4"
	_, err = tx.Exec(ctx, sql, key, record.Name, record.Value.Kind(), record.Value.String())
	if err != nil {
		rErr := tx.Rollback(ctx)
		if rErr != nil {
			return fmt.Errorf("db storage Push() -> Rollback() error: %w", err)
		}
		return fmt.Errorf("db storage Push() -> Exec() error: %w", err)
	}

	return tx.Commit(ctx)
}

// Push list of records to storage
func (d DatabaseStorage) PushList(ctx context.Context, data map[string]Record) error {
	batch := new(pgx.Batch)
	sql := "INSERT INTO metrics(id, name, kind, value) values ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET value = $4"
	for id, record := range data {
		batch.Queue(sql, id, record.Name, record.Value.Kind(), record.Value.String())
	}

	batchResp := d.Pool.SendBatch(ctx, batch)
	defer func() {
		if err := batchResp.Close(); err != nil {
			logging.LogErrorCtx(ctx, err, "failed to close batchResp")
		}
	}()

	for i := 0; i < len(data); i++ {
		if _, err := batchResp.Exec(); err != nil {
			return fmt.Errorf("db storage PushBatch() Exec error: %w", err)
		}
	}

	return nil
}

// Get a record from storage
func (d DatabaseStorage) Get(ctx context.Context, key string) (Record, error) {
	var (
		name   string
		kind   string
		value  float64
		record Record
		err    error
	)

	sql := "SELECT name, kind, value FROM metrics WHERE id=$1"
	err = d.Pool.QueryRow(ctx, sql, string(key)).Scan(&name, &kind, &value)

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

// Get list of records from storage
func (d DatabaseStorage) List(ctx context.Context) ([]Record, error) {
	rows, err := d.Pool.Query(ctx, "SELECT name, kind, value FROM metrics")
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

// Healthcheck
func (d DatabaseStorage) Ping(ctx context.Context) error {
	if err := d.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("db storage Ping() error: %w", err)
	}

	return nil
}

// Close storage pool
func (d DatabaseStorage) Close(ctx context.Context) error {
	d.Pool.Close()
	return nil
}
