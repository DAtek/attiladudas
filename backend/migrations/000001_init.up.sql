BEGIN;

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE gallery (
	id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR UNIQUE NOT NULL,
    slug VARCHAR UNIQUE NOT NULL,
    directory VARCHAR UNIQUE NOT NULL,
    "date" DATE NOT NULL,
	active BOOL NOT NULL DEFAULT FALSE,
    description VARCHAR
);

CREATE TABLE file (
	id SERIAL PRIMARY KEY NOT NULL,
	gallery_id INTEGER NOT NULL,
	filename VARCHAR NOT NULL,
    "rank" INTEGER NOT NULL,
    CONSTRAINT fk_file_gallery_id FOREIGN KEY(gallery_id) REFERENCES gallery(id) ON DELETE CASCADE,
    CONSTRAINT uk_file_gallery_id_filename UNIQUE(gallery_id, filename)
);

END;