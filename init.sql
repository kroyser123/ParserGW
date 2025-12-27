CREATE TABLE IF NOT EXISTS configs (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    version INTEGER,
    author TEXT,
    tags TEXT[]
);