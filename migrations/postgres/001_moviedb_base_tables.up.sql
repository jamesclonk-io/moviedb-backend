-- movie_dbdate
CREATE TABLE movie_dbdate (
    id              SERIAL PRIMARY KEY,
    date            TIMESTAMP NOT NULL DEFAULT CURRENT_DATE,
    description     TEXT NOT NULL
);

-- movie_genre
CREATE TABLE movie_genre (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE
);

-- movie_language
CREATE TABLE movie_language (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    country         TEXT NOT NULL,
    native_name     TEXT NOT NULL
);

-- movie_people
CREATE TABLE movie_people (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE
);

-- movie_movie
CREATE TABLE movie_movie (
    id              SERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    alttitle        TEXT,
    year            INTEGER,
    description     TEXT,
    format          TEXT,
    length          INTEGER,
    disk_region     TEXT,
    rating          INTEGER,
    disks           INTEGER,
    score           INTEGER,
    picture         TEXT,
    disk_type       TEXT
);