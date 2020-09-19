CREATE TABLE rv.user (
    id  BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
);

INSERT INTO rv.user DEFAULT VALUES;

ALTER TABLE rv.election
    ADD COLUMN user_id BIGINT NULL REFERENCES rv.user;

UPDATE rv.election
SET user_id = 1
WHERE user_id IS NULL;

ALTER TABLE rv.election
    ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE rv.vote
    ADD COLUMN user_id BIGINT NULL REFERENCES rv.user;

UPDATE rv.vote
SET user_id = 1
WHERE user_id IS NULL;

ALTER TABLE rv.vote
    ALTER COLUMN user_id SET NOT NULL;
