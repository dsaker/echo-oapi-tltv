CREATE EXTENSION IF NOT EXISTS pg_trgm ;

CREATE INDEX IF NOT EXISTS languages_language_trgm_idx ON languages USING gist (language gist_trgm_ops);