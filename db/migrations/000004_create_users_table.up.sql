CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
   id bigserial PRIMARY KEY,
   title_id bigint NOT NULL REFERENCES titles DEFAULT -1,
   name text UNIQUE NOT NULL,
   email citext UNIQUE NOT NULL,
   hashed_password varchar NOT NULL,
   flipped bool NOT NUll DEFAULT FALSE,
   og_language_id bigint NOT NULL REFERENCES languages DEFAULT -1,
   new_language_id bigint NOT NULL REFERENCES languages DEFAULT -1,
   created timestamp(0) NOT NULL DEFAULT NOW()
);
