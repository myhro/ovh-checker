BEGIN;

CREATE TABLE IF NOT EXISTS processor (
  id SERIAL PRIMARY KEY,
  brand TEXT,
  name TEXT
);

INSERT INTO processor (brand, name) VALUES
  ('AMD', 'Opteron'),
  ('Intel', 'Atom'),
  ('Intel', 'Core i3'),
  ('Intel', 'Core i5'),
  ('Intel', 'Core i7'),
  ('Intel', 'Xeon')
;

COMMIT;
