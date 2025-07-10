DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app_user') THEN
        CREATE ROLE app_user WITH LOGIN PASSWORD 'securepassword';
    END IF;
END
$$;

DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'app_db') THEN
        CREATE DATABASE app_db OWNER app_user;
    END IF;
END
$$;

\connect app_db

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,  -- Изменено с name на username
    password_hash TEXT NOT NULL     -- Добавлено вместо email
);

-- Даем права пользователю
GRANT ALL PRIVILEGES ON TABLE users TO app_user;
GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO app_user;