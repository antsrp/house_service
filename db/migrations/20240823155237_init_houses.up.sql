CREATE TABLE IF NOT EXISTS houses
(
    id SERIAL NOT NULL PRIMARY KEY,
    address VARCHAR(100) NOT NULL,
    year INTEGER,
    developer VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);