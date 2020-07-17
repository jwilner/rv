DROP TABLE rv.ballot;

ALTER TABLE rv.election
    ALTER COLUMN key DROP DEFAULT,
    ALTER COLUMN created_at DROP DEFAULT,
    ADD COLUMN ballot_key TEXT NOT NULL UNIQUE;

DROP FUNCTION generate_random_id;

DROP FUNCTION get_random_string;

CREATE TABLE rv.vote (
     id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
     election_id   BIGINT      NOT NULL REFERENCES rv.election (id),
     name          TEXT        NOT NULL,
     choices       TEXT[]      NOT NULL,
     created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
     UNIQUE (election_id, name)
)
