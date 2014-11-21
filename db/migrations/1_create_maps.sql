-- +migrate Up
CREATE TABLE maps (
  id serial PRIMARY KEY,
  name varchar(256) UNIQUE NOT NULL,
  structure text NOT NULL,
  created_at time,
  updated_at time
);

-- +migrate Down
DROP TABLE maps;
