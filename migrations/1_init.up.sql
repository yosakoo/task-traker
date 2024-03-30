 CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    name      TEXT NOT NULL,
    email     TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    UNIQUE (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

  
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    title TEXT NOT NULL,
    text TEXT,
    time TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);