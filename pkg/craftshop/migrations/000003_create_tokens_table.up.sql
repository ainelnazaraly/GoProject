CREATE TABLE IF NOT EXISTS tokens (
    plaintext text PRIMARY KEY,
    users_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

select * from tokens; 