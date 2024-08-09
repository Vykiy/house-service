CREATE TABLE IF NOT EXISTS flats (
    id SERIAL PRIMARY KEY,
    flat_number INTEGER NOT NULL,
    house_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    rooms INTEGER NOT NULL,
    status VARCHAR(255) NOT NULL,
    moderator_id UUID,
    CONSTRAINT unique_flat_number UNIQUE (house_id, flat_number),
    FOREIGN KEY (house_id) REFERENCES houses(id),
    FOREIGN KEY (moderator_id) REFERENCES users(id)
);