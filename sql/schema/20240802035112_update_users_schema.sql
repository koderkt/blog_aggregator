-- +goose Up
DROP TABLE users;

CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(20) NOT NULL
);

-- +goose Down
DROP TABLE users;

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(20) NOT NULL
);
