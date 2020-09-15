CREATE TYPE rv.election_flag AS ENUM (
    'public'
    );

ALTER TABLE rv.election
    ADD COLUMN close    TIMESTAMPTZ        NULL,
    ADD COLUMN close_tz TEXT               NULL,
    ADD COLUMN flags    rv.election_flag[] NOT NULL DEFAULT '{}';

ALTER TABLE rv.election
    RENAME COLUMN name TO question;
