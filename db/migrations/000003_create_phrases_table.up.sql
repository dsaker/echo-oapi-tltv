CREATE TABLE IF NOT EXISTS phrases (
    id bigserial PRIMARY KEY,
    title_id bigint NOT NULL REFERENCES titles ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS translates (
    phrase_id bigint NOT NULL REFERENCES phrases ON DELETE CASCADE,
    language_id smallint NOT NULL REFERENCES languages ON DELETE CASCADE,
    phrase text NOT NULL,
    phrase_hint text NOT NULL,
    PRIMARY KEY (language_id, phrase_id)
);
