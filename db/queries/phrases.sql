-- name: SelectUsersPhrasesByCorrect :many
SELECT phrase_id from users_phrases
WHERE user_id = $1 and title_id = $2 and language_id = $3
ORDER BY  phrase_correct, phrase_id
limit $4;

-- name: SelectUsersPhrasesByIds :one
SELECT * from users_phrases
WHERE user_id = $1 and language_id = $2 and phrase_id = $3;


-- name: SelectPhrasesFromTranslates :many
SELECT og.phrase_id, og.phrase, og.phrase_hint, new.phrase, new.phrase_hint
FROM
    (
        SELECT phrase_id, phrase, phrase_hint
        FROM translates t
        where t.language_id = $1
    ) og
    JOIN (
        SELECT phrase, phrase_hint, phrase_id
        FROM translates t
        WHERE t.language_id = $2
    ) new
ON og.phrase_id = new.phrase_id
WHERE og.phrase_id
IN (
    SELECT phrase_id from users_phrases
    WHERE user_id = $3 AND title_id = $4 AND language_id = $2
    ORDER BY  phrase_correct, phrase_id
    LIMIT $5
    );

-- name: SelectPhrasesFromTranslatesWithCorrect :many
SELECT og.phrase_id, og.phrase, og.phrase_hint, new.phrase, new.phrase_hint, up.phrase_correct
FROM
    (
        SELECT phrase_id, phrase, phrase_hint
        FROM translates t
        where t.language_id = $1
    ) og
    JOIN (
        SELECT phrase, phrase_hint, phrase_id
        FROM translates t
        WHERE t.language_id = $2
    ) new
ON og.phrase_id = new.phrase_id
    JOIN (
        SELECT phrase_id, phrase_correct
        FROM users_phrases u
        WHERE u.user_id = $3 AND u.title_id = $4 AND u.language_id = $2
    ) up
ON new.phrase_id = up.phrase_id
WHERE og.phrase_id
IN
    (
      SELECT phrase_id from users_phrases
      WHERE user_id = $3 AND title_id = $4 AND language_id = $2
      ORDER BY  phrase_correct, phrase_id
      LIMIT $5
    );

-- name: UpdateUsersPhrasesByThreeIds :one
UPDATE users_phrases
SET user_id = $1, title_id = $2, phrase_id = $3, language_id = $4, phrase_correct = $5
WHERE user_id = $1 AND phrase_id = $3 AND language_id = $4
RETURNING *;

-- name: InsertPhrases :one
INSERT INTO phrases (title_id)
VALUES ($1)
RETURNING *;

-- name: InsertTranslates :one
INSERT INTO translates (phrase_id, language_id, phrase, phrase_hint)
VALUES ($1, $2, $3, $4)
RETURNING *;