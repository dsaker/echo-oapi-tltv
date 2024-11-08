CREATE TABLE IF NOT EXISTS voices (
    id smallserial primary key not null ,
    language_id smallint references languages on delete cascade,
    language_codes text[] not null,
    ssml_gender text not null ,
    name text not null unique,
    natural_sample_rate_hertz smallint not null
);

ALTER TABLE voices
    ADD CONSTRAINT SsmlGenderCheck CHECK ( voices.ssml_gender SIMILAR TO 'FEMALE|MALE');
COMMENT ON CONSTRAINT SsmlGenderCheck ON voices IS 'One of FEMALE|MALE';