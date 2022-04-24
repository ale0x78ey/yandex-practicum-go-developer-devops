CREATE TABLE counter_metrics (
  id         text NOT NULL,
  created_at timestamp DEFAULT (now() at time zone 'UTC') NOT NULL,
  value      bigint NOT NULL
);
