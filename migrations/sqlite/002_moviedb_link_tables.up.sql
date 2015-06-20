-- movie_link_actor
CREATE TABLE IF NOT EXISTS `movie_link_actor` (
	`movie_id`			integer NOT NULL,
	`person_id`			integer NOT NULL,
	FOREIGN KEY(`movie_id`) REFERENCES [movie_movie] ( [id] ) ON DELETE CASCADE,
	FOREIGN KEY(`person_id`) REFERENCES [movie_people] ( [id] ) ON DELETE CASCADE
);

-- movie_link_director
CREATE TABLE IF NOT EXISTS `movie_link_director` (
	`movie_id`			integer NOT NULL,
	`person_id`			integer NOT NULL,
	FOREIGN KEY(`movie_id`) REFERENCES [movie_movie] ( [id] ) ON DELETE CASCADE,
	FOREIGN KEY(`person_id`) REFERENCES [movie_people] ( [id] ) ON DELETE CASCADE
);

-- movie_link_genre
CREATE TABLE IF NOT EXISTS `movie_link_genre` (
	`movie_id`			integer NOT NULL,
	`genre_id`			integer NOT NULL,
	FOREIGN KEY(`movie_id`) REFERENCES [movie_movie] ( [id] ) ON DELETE CASCADE,
	FOREIGN KEY(`genre_id`) REFERENCES [movie_genre] ( [id] ) ON DELETE CASCADE
);

-- movie_link_language
CREATE TABLE IF NOT EXISTS `movie_link_language` (
	`movie_id`			integer NOT NULL,
	`language_id`		integer NOT NULL,
	FOREIGN KEY(`movie_id`) REFERENCES [movie_movie] ( [id] ) ON DELETE CASCADE,
	FOREIGN KEY(`language_id`) REFERENCES [movie_language] ( [id] ) ON DELETE CASCADE
);
