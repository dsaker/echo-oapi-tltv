CREATE TABLE IF NOT EXISTS languages (
    id smallserial PRIMARY KEY ,
    language text NOT NULL,
    tag text NOT NULL
);

INSERT INTO languages (id, language, tag) VALUES (-1, 'Not a Language', 'NaL');