-- +migrate Up
ALTER TABLE maps
  ADD sort_by varchar(256);

UPDATE maps
  SET sort_by = ''
  WHERE sort_by IS NULL;


-- +migrate Down
ALTER TABLE maps
  DROP COLUMN sort_by;
