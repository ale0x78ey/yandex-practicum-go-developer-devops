CREATE TABLE gauge_metrics (
  id         text NOT NULL,
  created_at timestamp DEFAULT (now() at time zone 'UTC') NOT NULL,
  value      double precision NOT NULL
);
