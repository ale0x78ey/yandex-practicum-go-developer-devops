package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type Config struct {
	DSN           string `env:"DATABASE_DSN"`
	MigrationsURL string
}

type MetricStorage struct {
	config              Config
	db                  *sql.DB
	gaugeSaveStmt       *sql.Stmt
	gaugeIncrStmt       *sql.Stmt
	gaugeLoadStmt       *sql.Stmt
	gaugeLoadListStmt   *sql.Stmt
	counterSaveStmt     *sql.Stmt
	counterIncrStmt     *sql.Stmt
	counterLoadStmt     *sql.Stmt
	counterLoadListStmt *sql.Stmt
}

func NewMetricStorage(ctx context.Context, config Config) (*MetricStorage, error) {
	db, err := sql.Open("pgx", config.DSN)
	if err != nil {
		return nil, err
	}

	storage := &MetricStorage{
		config: config,
		db:     db,
	}

	if err := storage.Migrate(); err != nil {
		db.Close()
		return nil, err
	}

	if err := storage.PrepareStatements(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *MetricStorage) Migrate() error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	cfg := &pgx.Config{}
	driver, err := pgx.WithInstance(s.db, cfg)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(s.config.MigrationsURL, cfg.DatabaseName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *MetricStorage) PrepareStatements(ctx context.Context) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	gaugeSaveStmt, err := s.db.PrepareContext(
		ctx,
		`
INSERT INTO gauge_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = $2`,
	)
	if err != nil {
		return err
	}
	s.gaugeSaveStmt = gaugeSaveStmt

	gaugeIncrStmt, err := s.db.PrepareContext(
		ctx,
		`
INSERT INTO gauge_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value + $2`,
	)
	if err != nil {
		return err
	}
	s.gaugeIncrStmt = gaugeIncrStmt

	gaugeLoadStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT value FROM gauge_metrics WHERE id = $1",
	)
	if err != nil {
		return err
	}
	s.gaugeLoadStmt = gaugeLoadStmt

	gaugeLoadListStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT id, value FROM gauge_metrics",
	)
	if err != nil {
		return err
	}
	s.gaugeLoadListStmt = gaugeLoadListStmt

	counterSaveStmt, err := s.db.PrepareContext(
		ctx,
		`
INSERT INTO counter_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = $2`,
	)
	if err != nil {
		return err
	}
	s.counterSaveStmt = counterSaveStmt

	counterIncrStmt, err := s.db.PrepareContext(
		ctx,
		`
INSERT INTO counter_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value + $2`,
	)
	if err != nil {
		return err
	}
	s.counterIncrStmt = counterIncrStmt

	counterLoadStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT value FROM counter_metrics WHERE id = $1",
	)
	if err != nil {
		return err
	}
	s.counterLoadStmt = counterLoadStmt

	counterLoadListStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT id, value FROM counter_metrics",
	)
	if err != nil {
		return err
	}
	s.counterLoadListStmt = counterLoadListStmt

	return nil
}

func (s *MetricStorage) SaveMetric(ctx context.Context, metric model.Metric) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	if err := metric.Validate(); err != nil {
		return err
	}

	switch metric.MType {
	case model.MetricTypeGauge:
		if _, err := s.gaugeSaveStmt.ExecContext(ctx, metric.ID, *metric.Value); err != nil {
			return err
		}
	case model.MetricTypeCounter:
		if _, err := s.counterSaveStmt.ExecContext(ctx, metric.ID, *metric.Delta); err != nil {
			return err
		}
	}

	return nil
}

func (s *MetricStorage) IncrMetric(ctx context.Context, metric model.Metric) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	if err := metric.Validate(); err != nil {
		return err
	}

	switch metric.MType {
	case model.MetricTypeGauge:
		if _, err := s.gaugeIncrStmt.ExecContext(ctx, metric.ID, *metric.Value); err != nil {
			return err
		}
	case model.MetricTypeCounter:
		if _, err := s.counterIncrStmt.ExecContext(ctx, metric.ID, *metric.Delta); err != nil {
			return err
		}
	}

	return nil
}

func (s *MetricStorage) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (*model.Metric, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not opened")
	}

	if metricName == "" {
		return nil, errors.New("invalid empty metricName")
	}

	metric := &model.Metric{ID: metricName}

	var err error

	switch metricType {
	case model.MetricTypeGauge:
		row := s.gaugeLoadStmt.QueryRowContext(ctx, metricName)
		err = row.Scan(&metric.Value)
	case model.MetricTypeCounter:
		row := s.counterLoadStmt.QueryRowContext(ctx, metricName)
		err = row.Scan(&metric.Delta)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *MetricStorage) loadGaugeMetricList(
	ctx context.Context,
	tx *sql.Tx,
) ([]model.Metric, error) {
	txStmt := tx.StmtContext(ctx, s.gaugeLoadStmt)
	rows, err := txStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0, 50)

	for rows.Next() {
		metric := model.Metric{}
		if err := rows.Scan(&metric.ID, metric.Value); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *MetricStorage) loadCounterMetricList(
	ctx context.Context,
	tx *sql.Tx,
) ([]model.Metric, error) {
	txStmt := tx.StmtContext(ctx, s.counterLoadStmt)
	rows, err := txStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0, 50)

	for rows.Next() {
		metric := model.Metric{}
		if err := rows.Scan(&metric.ID, metric.Delta); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *MetricStorage) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not opened")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	gaugeMetrics, err := s.loadGaugeMetricList(ctx, tx)
	if err != nil {
		return nil, err
	}

	counterMetrics, err := s.loadCounterMetricList(ctx, tx)
	if err != nil {
		return nil, err
	}

	metrics := make([]model.Metric, len(gaugeMetrics)+len(counterMetrics))
	_ = copy(metrics, gaugeMetrics)
	_ = copy(metrics[len(gaugeMetrics):], counterMetrics)

	return metrics, tx.Commit()
}

func (s *MetricStorage) Flush(ctx context.Context) error {
	return nil
}

func (s *MetricStorage) Close() {
	for _, stmt := range []*sql.Stmt{
		s.gaugeSaveStmt,
		s.gaugeIncrStmt,
		s.gaugeLoadStmt,
		s.gaugeLoadListStmt,
		s.counterSaveStmt,
		s.counterIncrStmt,
		s.counterLoadStmt,
		s.counterLoadListStmt,
	} {
		if stmt != nil {
			stmt.Close()
		}
	}

	s.db.Close()
}

func (s *MetricStorage) Validate(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
