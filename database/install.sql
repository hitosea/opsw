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

CREATE TABLE IF NOT EXISTS servers (
    id integer NOT NULL PRIMARY KEY {{.INCREMENT}},
    ip VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    port VARCHAR(100) NOT NULL,
    remark VARCHAR(255) NOT NULL,
    state VARCHAR(50) NOT NULL,
    token VARCHAR(100) NOT NULL,
    systems TEXT NOT NULL,
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);

CREATE TABLE IF NOT EXISTS server_users (
    id integer NOT NULL PRIMARY KEY {{.INCREMENT}},
    user_id integer NOT NULL,
    server_id integer NOT NULL,
    owner_id integer NOT NULL,
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);

CREATE TABLE IF NOT EXISTS server_infos (
    id integer NOT NULL PRIMARY KEY {{.INCREMENT}},
    server_id integer NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    os VARCHAR(255) NOT NULL,
    platform VARCHAR(255) NOT NULL,
    platform_family VARCHAR(255) NOT NULL,
    platform_version VARCHAR(255) NOT NULL,
    kernel_arch VARCHAR(255) NOT NULL,
    kernel_version VARCHAR(255) NOT NULL,
    virtualization_system VARCHAR(255) NOT NULL,
    cpu_cores integer NOT NULL,
    cpu_logical_cores integer NOT NULL,
    cpu_model_name VARCHAR(255) NOT NULL,
    current_info TEXT NOT NULL,
    version VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(255) NOT NULL,
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);