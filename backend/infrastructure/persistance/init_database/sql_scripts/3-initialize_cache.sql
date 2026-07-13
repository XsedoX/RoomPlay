CREATE UNLOGGED TABLE IF NOT EXISTS cache
(
  id SERIAL PRIMARY KEY,
  key TEXT  NOT NULL,
  token TEXT NOT NULL DEFAULT '',
  value JSONB,
  created_at_utc TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (key, token)
);

CREATE INDEX IF NOT EXISTS cache_idx ON cache(id) INCLUDE (value);
CREATE INDEX IF NOT EXISTS cache_key_trgm_idx ON cache USING gin (key gin_trgm_ops);
