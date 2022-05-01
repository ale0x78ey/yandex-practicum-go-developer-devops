package db

import (
	"context"
	"errors"
)

func (s *MetricStorage) PrepareStatements(ctx context.Context) error {
	if s.db == nil {
		return errors.New("database connection is not opened")
	}

	if err := s.prepareGaugeSaveStmt(ctx); err != nil {
		return err
	}

	if err := s.prepareGaugeIncrStmt(ctx); err != nil {
		return err
	}

	if err := s.prepareGaugeLoadStmt(ctx); err != nil {
		return err
	}

	if err := s.prepareGaugeLoadListStmt(ctx); err != nil {
		return err
	}

	if err := s.counterGaugeSaveStmt(ctx); err != nil {
		return err
	}

	if err := s.counterGaugeIncrStmt(ctx); err != nil {
		return err
	}

	if err := s.counterGaugeLoadStmt(ctx); err != nil {
		return err
	}

	if err := s.counterGaugeLoadListStmt(ctx); err != nil {
		return err
	}

	return nil
}

func (s *MetricStorage) prepareGaugeSaveStmt(ctx context.Context) error {
	expr := `
INSERT INTO gauge_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = $2`

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.gaugeSaveStmt = stmt
	return nil
}

func (s *MetricStorage) prepareGaugeIncrStmt(ctx context.Context) error {
	expr := `
INSERT INTO gauge_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value + $2`

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.gaugeIncrStmt = stmt
	return nil
}

func (s *MetricStorage) prepareGaugeLoadStmt(ctx context.Context) error {
	expr := "SELECT value FROM gauge_metrics WHERE id = $1"
	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.gaugeLoadStmt = stmt
	return nil
}

func (s *MetricStorage) prepareGaugeLoadListStmt(ctx context.Context) error {
	expr := "SELECT id, value FROM gauge_metrics"
	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.gaugeLoadListStmt = stmt
	return nil
}

func (s *MetricStorage) counterGaugeSaveStmt(ctx context.Context) error {
	expr := `
INSERT INTO counter_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = $2`

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.counterSaveStmt = stmt
	return nil
}

func (s *MetricStorage) counterGaugeIncrStmt(ctx context.Context) error {
	expr := `
INSERT INTO counter_metrics (id, value)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value + $2`

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.counterIncrStmt = stmt
	return nil
}

func (s *MetricStorage) counterGaugeLoadStmt(ctx context.Context) error {
	expr := "SELECT value FROM counter_metrics WHERE id = $1"

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.counterLoadStmt = stmt
	return nil
}

func (s *MetricStorage) counterGaugeLoadListStmt(ctx context.Context) error {
	expr := "SELECT id, value FROM counter_metrics"

	stmt, err := s.db.PrepareContext(ctx, expr)
	if err != nil {
		return err
	}
	s.counterLoadListStmt = stmt
	return nil
}
