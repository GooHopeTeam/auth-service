CREATE TABLE IF NOT EXISTS "user"
(
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(320) NOT NULL UNIQUE,
    hashed_password VARCHAR(256) NOT NULL,
    created_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS token
(
    user_id INTEGER REFERENCES "user",
    token   VARCHAR(128) NOT NULL UNIQUE
);
