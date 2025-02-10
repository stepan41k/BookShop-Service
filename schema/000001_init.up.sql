CREATE TABLE
    authors (
        id SERIAL NOT NULL UNIQUE,
        author VARCHAR(50) NOT NULL
    );

CREATE TABLE
    genres (
        id SERIAL NOT NULL UNIQUE,
        genre VARCHAR(50) NOT NULL UNIQUE
    );

CREATE TABLE
    events (
        id SERIAL NOT NULL UNIQUE,
        event_type TEXT NOT NULL,
        payload TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'done')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    books (
        id SERIAL NOT NULL UNIQUE,
        name VARCHAR(50),
        author_id INT REFERENCES authors (id) ON DELETE CASCADE NOT NULL,
        genre_id INT REFERENCES genres (id) ON DELETE CASCADE NOT NULL,
        price INT NOT NULL
    );
