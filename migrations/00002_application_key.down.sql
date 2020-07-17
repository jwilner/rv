DROP TABLE rv.vote;

-- random ID generation logic for postgres from https://www.depesz.com/2017/02/06/generate-short-random-textual-ids/
-- Ideally, we'd be able to install this as a PG extension but I don't want to figure out how to do that with heroku.
CREATE FUNCTION get_random_string(IN string_length integer,
                                  IN possible_chars TEXT DEFAULT '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz') RETURNS text
    LANGUAGE plpgsql
AS
$$
DECLARE
    output TEXT = '';
    i      INT4;
    pos    INT4;
BEGIN
    FOR i IN 1..string_length
        LOOP
            pos := 1 + cast(random() * (length(possible_chars) - 1) as INT4);
            output := output || substr(possible_chars, pos, 1);
        END LOOP;
    RETURN output;
END;
$$;

CREATE OR REPLACE FUNCTION generate_random_id(IN table_schema TEXT,
                                              IN table_name TEXT,
                                              IN column_name TEXT,
                                              IN string_length INTEGER) returns text
    LANGUAGE plpgsql
AS
$$
DECLARE
    v_random_id   text;
    v_temp        text;
    v_length      int4 := string_length;
    v_sql         text;
    v_advisory_1  int4 := hashtext(format('%I:%I:%I', table_schema, table_name, column_name));
    v_advisory_2  int4;
    v_advisory_ok bool;
BEGIN
    v_sql := format('SELECT %I FROM %I.%I WHERE %I = $1', column_name, table_schema, table_name, column_name);
    LOOP
        v_random_id := get_random_string(v_length, '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz');
        v_advisory_2 := hashtext(v_random_id);
        v_advisory_ok := pg_try_advisory_xact_lock(v_advisory_1, v_advisory_2);
        IF v_advisory_ok THEN
            EXECUTE v_sql INTO v_temp USING v_random_id;
            exit when v_temp is null;
        END IF;
        v_length := v_length + 1;
    END LOOP;
    return v_random_id;
END;
$$
    STRICT
    SET search_path TO rv;
;

ALTER TABLE rv.ballot
    ALTER COLUMN key SET DEFAULT generate_random_id('rv', 'ballot', 'key', 8),
    ALTER COLUMN created_at SET DEFAULT now(),
    ADD COLUMN name TEXT NOT NULL DEFAULT '';

ALTER TABLE rv.election
    ALTER COLUMN key SET DEFAULT generate_random_id('rv', 'election', 'key', 8),
    ALTER COLUMN created_at SET DEFAULT now(),
    ADD COLUMN choices TEXT[] NOT NULL DEFAULT '{}'::TEXT[];

