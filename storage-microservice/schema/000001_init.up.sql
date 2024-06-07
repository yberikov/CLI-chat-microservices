CREATE TABLE messages (
    id serial not null unique,
    text TEXT  not null,
    author VARCHAR(25) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
