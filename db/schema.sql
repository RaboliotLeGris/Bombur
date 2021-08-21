-- SCHEMA

SET timezone = 'Europe/Paris';

CREATE TABLE meta_info (
    version INT PRIMARY KEY
);

CREATE TABLE link (
    id  INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    link TEXT NOT NULL,
    slug TEXT NOT NULL,
    expire TIMESTAMPTZ
);

CREATE UNIQUE INDEX slug_idx ON link (slug);

-- initial data

INSERT INTO meta_info VALUES (1);