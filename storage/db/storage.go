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
	config Config

	db *sql.DB

	gaugeSaveStmt     *sql.Stmt
	gaugeIncrStmt     *sql.Stmt
	gaugeLoadStmt     *sql.Stmt
	gaugeLoadListStmt *sql.Stmt

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
		storage.Close()
		return nil, err
	}

	if err := storage.PrepareStatements(ctx); err != nil {
		storage.Close()
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
		if _, err := s.counterSaveStmt.ExecContext(ctx, metric.ID, *metric.Value); err != nil {
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
		if _, err := s.counterIncrStmt.ExecContext(ctx, metric.ID, *metric.Value); err != nil {
			return err
		}
	}

	return nil
}

func (s *MetricStorage) LoadMetric(
	ctx context.Context,
	metric model.Metric,
) (*model.Metric, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not opened")
	}

	if err := metric.ID.Validate(); err != nil {
		return nil, err
	}

	if err := metric.MType.Validate(); err != nil {
		return nil, err
	}

	var err error
	switch metric.MType {
	case model.MetricTypeGauge:
		row := s.gaugeLoadStmt.QueryRowContext(ctx, metric.ID)
		value := model.Gauge(0)
		metric.Value = &value
		err = row.Scan(metric.Value)

	case model.MetricTypeCounter:
		row := s.counterLoadStmt.QueryRowContext(ctx, metric.ID)
		delta := model.Counter(0)
		metric.Delta = &delta
		err = row.Scan(metric.Delta)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (s *MetricStorage) loadGaugeMetricList(
	ctx context.Context,
	tx *sql.Tx,
) ([]model.Metric, error) {
	txStmt := tx.StmtContext(ctx, s.gaugeLoadListStmt)
	rows, err := txStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0, 50)

	for rows.Next() {
		metric := model.MetricFromGauge("", 0)
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
	txStmt := tx.StmtContext(ctx, s.counterLoadListStmt)
	rows, err := txStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0, 50)

	for rows.Next() {
		metric := model.MetricFromCounter("", 0)
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

func (s *MetricStorage) SaveMetricList(ctx context.Context, metrics []model.Metric) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txGaugeSaveStmt := tx.StmtContext(ctx, s.gaugeSaveStmt)
	txCounterSaveStmt := tx.StmtContext(ctx, s.counterSaveStmt)

	for _, m := range metrics {
		switch m.MType {
		case model.MetricTypeGauge:
			if _, err := txGaugeSaveStmt.ExecContext(ctx, m.ID, *m.Value); err != nil {
				return err
			}

		case model.MetricTypeCounter:
			if _, err := txCounterSaveStmt.ExecContext(ctx, m.ID, *m.Delta); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *MetricStorage) IncrMetricList(ctx context.Context, metrics []model.Metric) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txGaugeIncrStmt := tx.StmtContext(ctx, s.gaugeIncrStmt)
	txCounterIncrStmt := tx.StmtContext(ctx, s.counterIncrStmt)

	for _, m := range metrics {
		switch m.MType {
		case model.MetricTypeGauge:
			if _, err := txGaugeIncrStmt.ExecContext(ctx, m.ID, *m.Value); err != nil {
				return err
			}

		case model.MetricTypeCounter:
			if _, err := txCounterIncrStmt.ExecContext(ctx, m.ID, *m.Delta); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
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

	// if s.db == nil {
	// 	return errors.New("database connection is not opened")
	// }

	// tx, err := s.db.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer tx.Rollback()

	// txGaugeSaveStmt := tx.StmtContext(ctx, s.gaugeSaveStmt)
	// txCounterSaveStmt := tx.StmtContext(ctx, s.counterSaveStmt)

	// for _, m := range s.mCache.GetAll() {
	// 	switch m.MType {
	// 	case model.MetricTypeGauge:
	// 		if _, err := txGaugeSaveStmt.ExecContext(ctx, m.ID, *m.Value); err != nil {
	// 			return err
	// 		}
	// 	case model.MetricTypeCounter:
	// 		if _, err := txCounterSaveStmt.ExecContext(ctx, m.ID, *m.Delta); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// return tx.Commit()
}

func (s *MetricStorage) Close() {
	for _, stmt := range []*sql.Stmt{
		s.gaugeSaveStmt,
		s.gaugeLoadStmt,
		s.gaugeLoadListStmt,
		s.counterSaveStmt,
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
