CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INT,
    status TEXT NOT NULL DEFAULT 'active',

    CONSTRAINT fk_user_refresh_token FOREIGN KEY (user_id) REFERENCES users(id)
);
