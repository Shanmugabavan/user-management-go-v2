CREATE TABLE users (
user_id     UUID PRIMARY KEY,
first_name  TEXT NOT NULL,
last_name   TEXT NOT NULL,
email       TEXT NOT NULL UNIQUE,
phone       TEXT NOT NULL,
age         INT NOT NULL CHECK (age >= 0),
status      INT NOT NULL
);