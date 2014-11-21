-- +migrate Up
CREATE TABLE map_items (
  id serial PRIMARY KEY,
  structure_data text NOT NULL,
  created_at time,
  updated_at time
);

-- +migrate Down
DROP TABLE map_items;
