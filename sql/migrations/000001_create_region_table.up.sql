BEGIN;

CREATE TABLE IF NOT EXISTS region (
  id SERIAL PRIMARY KEY,
  code TEXT,
  name TEXT
);

INSERT INTO region (code, name) VALUES
  ('apac', 'Asia Pacific'),
  ('europe', 'Europe'),
  ('northAmerica', 'North America')
;

COMMIT;
