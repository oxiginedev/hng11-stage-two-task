-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
    id CHAR(26) PRIMARY KEY,

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS organisations(
    id CHAR(26) PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS organisation_user(
    id CHAR(26) PRIMARY KEY,

    organisation_id CHAR(26) NOT NULL REFERENCES organisations (id),
    user_id CHAR(26) NOT NULL REFERENCES users (id),

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS organisations CASCADE;
DROP TABLE IF EXISTS organisation_user;
