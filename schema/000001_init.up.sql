CREATE TABLE
    book (
        id SERIAL NOT NULL UNIQUE,
        name VARCHAR(50),
        author_id INT REFERENCES author (id) ON DELETE CASCADE NOT NULL,
        genre_id INT REFERENCES genre (id) ON DELETE CASCADE NOT NULL,
        price INT NOT NULL
    );

CREATE TABLE
    genre (
        id SERIAL NOT NULL UNIQUE,
        genre VARCHAR(50) NOT NULL UNIQUE
    );

CREATE TABLE
    author (
        id SERIAL NOT NULL UNIQUE,
        author VARCHAR(50) NOT NULL
    );