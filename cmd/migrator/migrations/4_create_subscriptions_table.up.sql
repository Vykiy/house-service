CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT unique_subscription UNIQUE (house_id, email),
    FOREIGN KEY (house_id) REFERENCES houses(id)
);