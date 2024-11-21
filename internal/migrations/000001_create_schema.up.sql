CREATE TABLE author (
    author_id   SERIAL PRIMARY KEY,
    name        VARCHAR(30) NOT NULL
);

CREATE TABLE song (
    song_id     SERIAL PRIMARY KEY,
    name        VARCHAR(30) NOT NULL,
    author      INT,
    text        TEXT,
    release     DATE,
    link        VARCHAR(300),
    FOREIGN KEY (author) REFERENCES author (author_id) ON DELETE SET NULL
);