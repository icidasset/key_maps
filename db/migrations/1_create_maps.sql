-- +migrate Up
CREATE TABLE maps (
  id serial PRIMARY KEY,
  name varchar(256) UNIQUE NOT NULL,
  slug varchar(256) UNIQUE NOT NULL,
  structure text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

-- +migrate Down
DROP TABLE maps;
