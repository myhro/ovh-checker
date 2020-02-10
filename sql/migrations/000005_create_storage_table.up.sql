BEGIN;

CREATE TABLE IF NOT EXISTS storage (
  id SERIAL PRIMARY KEY,
  code TEXT,
  size INTEGER,
  quantity INTEGER,
  type TEXT
);

INSERT INTO storage (code, size, quantity, type) VALUES
  ('500GB', 500, 1, 'HDD'),
  ('1TB', 1024, 1, 'HDD'),
  ('2TB', 2048, 1, 'HDD'),
  ('2x2TB', 2048, 2, 'HDD'),
  ('2x240GB', 240, 2, 'SSD')
;

COMMIT;
