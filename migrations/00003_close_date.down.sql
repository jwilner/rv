ALTER TABLE rv.election
    RENAME COLUMN name TO question;

ALTER TABLE rv.election
    DROP COLUMN close,
    DROP COLUMN close_tz,
    DROP COLUMN flags;

DROP TYPE rv.election_flag;