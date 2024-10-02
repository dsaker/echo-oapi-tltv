CREATE TABLE IF NOT EXISTS users_phrases (
     user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
     title_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
     phrase_id BIGINT NOT NULL REFERENCES phrases ON DELETE CASCADE,
     language_id INT NOT NULL REFERENCES phrases ON DELETE CASCADE,
     phrase_correct INT NOT NULL DEFAULT 0,
     PRIMARY KEY (user_id, phrase_id, language_id)
    );
