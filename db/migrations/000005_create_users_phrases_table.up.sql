CREATE TABLE IF NOT EXISTS users_phrases (
     user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
     phrase_id bigint NOT NULL REFERENCES phrases ON DELETE CASCADE,
     title_id bigint NOT NULL REFERENCES titles ON DELETE CASCADE,
     phrase_correct bigint,
     flipped_correct bigint,
     PRIMARY KEY (user_id, phrase_id, title_id)
    );
