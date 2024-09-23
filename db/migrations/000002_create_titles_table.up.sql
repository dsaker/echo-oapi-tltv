CREATE TABLE IF NOT EXISTS titles (
    id bigserial PRIMARY KEY,
    title text NOT NULL UNIQUE,
    num_subs int NOT NULL,
    language_id bigint NOT NULL REFERENCES languages default -1,
    og_language_id bigint NOT NULL REFERENCES languages default -1
);

INSERT INTO titles (id, title, num_subs) VALUES (-1, 'Not a Movie', 0);
