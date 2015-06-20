-- movie_dbdate
CREATE TABLE IF NOT EXISTS `movie_dbdate` (
    `id`            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `date`          datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
    `description`   text NOT NULL
);

-- movie_genre
CREATE TABLE IF NOT EXISTS `movie_genre` (
    `id`            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name`          text NOT NULL UNIQUE
);

-- movie_language
CREATE TABLE IF NOT EXISTS `movie_language` (
    `id`            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name`          text NOT NULL UNIQUE,
    `country`       text NOT NULL,
    `native_name`   text NOT NULL
);

-- movie_people
CREATE TABLE IF NOT EXISTS `movie_people` (
    `id`            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name`          text NOT NULL UNIQUE
);

-- movie_movie
CREATE TABLE IF NOT EXISTS `movie_movie` (
    `id`            integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `title`         text NOT NULL,
    `alttitle`      text,
    `year`          integer,
    `description`   text,
    `format`        text,
    `length`        integer,
    `disk_region`   text,
    `rating`        integer,
    `disks`         integer,
    `score`         integer,
    `picture`       text,
    `disk_type`     text
);