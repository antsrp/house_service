CREATE TABLE IF NOT EXISTS subscribers(
    id SERIAL NOT NULL PRIMARY KEY,
    email VARCHAR(50) NOT NULL,
    house_id INTEGER REFERENCES houses(id) ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sub ON subscribers(email, house_id);