CREATE TABLE gauge_metrics (
  id    text NOT NULL,
  value double precision NOT NULL,
  UNIQUE (id)
);
