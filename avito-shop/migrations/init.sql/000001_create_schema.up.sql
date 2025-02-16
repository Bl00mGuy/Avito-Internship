CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username TEXT UNIQUE,
    password TEXT,
    coins    INTEGER
);

CREATE TABLE IF NOT EXISTS purchases
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id),
    item       TEXT,
    price      INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS coin_transfers
(
    id         SERIAL PRIMARY KEY,
    from_user  INTEGER REFERENCES users (id),
    to_user    INTEGER REFERENCES users (id),
    amount     INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);