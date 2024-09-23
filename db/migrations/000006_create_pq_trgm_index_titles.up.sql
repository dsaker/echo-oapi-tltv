CREATE EXTENSION IF NOT EXISTS pg_trgm ;

CREATE INDEX IF NOT EXISTS titles_title_trgm_idx ON titles USING gist (title gist_trgm_ops);
