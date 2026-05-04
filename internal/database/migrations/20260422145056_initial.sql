-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    formatted_name TEXT,
    email TEXT NOT NULL UNIQUE COLLATE NOCASE,
    avatar TEXT,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS user_access (
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    issued_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    PRIMARY KEY (user_id, token_hash),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    category TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    onclick_path TEXT,
    image TEXT,
    level INTEGER DEFAULT 0, -- 0: info | 1: warning | 2: error | 3: success
    dismiss_timeout INTEGER, -- null or 0: default | -1: will not auto dismiss
    is_read BOOLEAN DEFAULT false,
    metadata JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS apps (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    project_name TEXT UNIQUE, -- the one that Compose prefixes for each stack of containers
    origin TEXT, -- url pointing to where this was taken
    compose_yaml TEXT, -- snapshot
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS app_volumes (
    id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL,
    name TEXT,
    driver TEXT,
    is_exposed BOOLEAN DEFAULT false, -- whether is exposed on Kaeru's file explorer
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_id) references apps (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS app_services (
    app_id TEXT NOT NULL,
    name TEXT NOT NULL,
    image TEXT,
    last_known_status TEXT,
    container_id TEXT, -- snapshot
    PRIMARY KEY (app_id, name),
    FOREIGN KEY (app_id) REFERENCES apps (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS app_service_volumes (
    service_id TEXT NOT NULL,
    service_name TEXT NOT NULL,
    volume_id TEXT NOT NULL,
    mount_path TEXT NOT NULL,
    readonly BOOLEAN DEFAULT false,
    type TEXT DEFAULT "volume", -- volume | bind
    bind_source TEXT, -- if type == bind
    PRIMARY KEY (service_id, volume_id),
    FOREIGN KEY (service_id, service_name) REFERENCES app_services (app_id, name) ON DELETE CASCADE,
    FOREIGN KEY (volume_id) REFERENCES app_volumes (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS app_configs (
    id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    is_sensitive BOOLEAN DEFAULT false, -- will hash value if true
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    UNIQUE (app_id, key),
    FOREIGN KEY (app_id) REFERENCES apps (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    status TEXT NOT NULL, -- queued | running | success | error
    progress INTEGER DEFAULT 0,
    logs TEXT,
    app_id TEXT, -- most jobs are probably app related, but not all, so can be nullable
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME, -- if status == success | error
    FOREIGN KEY (app_id) REFERENCES apps (id) -- better set up a CRON to clean up old jobs over time
);

-- +goose Down
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS user_access;

DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS apps;

DROP TABLE IF EXISTS app_volumes;

DROP TABLE IF EXISTS app_services;

DROP TABLE IF EXISTS app_service_volumes;

DROP TABLE IF EXISTS app_configs;

DROP TABLE IF EXISTS jobs;