CREATE TABLE IF NOT EXISTS phrases (
    id bigserial PRIMARY KEY,
    title_id bigint NOT NULL REFERENCES titles ON DELETE CASCADE,
    phrase text NOT NULL,
    translates text NOT NULL,
    phrase_hint text NOT NULL,
    translates_hint text NOT NULL
)
