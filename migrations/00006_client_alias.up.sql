CREATE TABLE rv.client
(
    id     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name   TEXT  NOT NULL UNIQUE,
    secret BYTEA NOT NULL
);

CREATE TABLE rv.alias
(
    id        BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id   BIGINT NOT NULL REFERENCES rv.user (id),
    client_id BIGINT NOT NULL REFERENCES rv.client (id),
    alias     TEXT   NOT NULL,

    UNIQUE (client_id, alias)
);
