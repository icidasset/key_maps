-- +migrate Up
CREATE TABLE map_items (
  id serial PRIMARY KEY,
  structure_data text NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  map_id integer NOT NULL
);

CREATE INDEX map_items_map_id_index ON map_items(map_id)

-- +migrate Down
DROP TABLE map_items;
