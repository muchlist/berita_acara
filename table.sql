CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    name VARCHAR (100) NOT NULL,
    email VARCHAR ( 255 ) UNIQUE NOT NULL,
    password VARCHAR (100) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS roles(
    role_name VARCHAR (20) PRIMARY KEY,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions(
    id SERIAL PRIMARY KEY,
    roles_name VARCHAR(20) REFERENCES roles(role_name) ON DELETE CASCADE,
    permission_name VARCHAR (20) UNIQUE NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_roles(
    id SERIAL PRIMARY KEY,
    users_id INT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    roles_name VARCHAR (20) REFERENCES roles(role_name) ON DELETE CASCADE ON UPDATE CASCADE
);
