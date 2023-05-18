CREATE TABLE IF NOT EXISTS users (
    id integer NOT NULL PRIMARY KEY {{.INCREMENT}},
    email VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    encrypt VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    token VARCHAR(100) NOT NULL,
    avatar VARCHAR(255) NOT NULL,
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);