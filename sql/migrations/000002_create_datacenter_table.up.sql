BEGIN;

CREATE TABLE IF NOT EXISTS datacenter (
  id SERIAL PRIMARY KEY,
  code TEXT,
  name TEXT,
  region_id INTEGER REFERENCES region(id)
);

INSERT INTO datacenter (code, name, region_id)
  SELECT data.code, data.name, region.id
    FROM (
      VALUES
        ('bhs', 'Beauharnois', 'northAmerica'),
        ('default', 'Default', 'europe'),
        ('fra', 'Frankfurt', 'europe'),
        ('gra', 'Gravelines', 'europe'),
        ('hil', 'Hillsboro', 'northAmerica'),
        ('lon', 'London', 'europe'),
        ('rbx', 'Roubaix', 'europe'),
        ('sbg', 'Strasbourg', 'europe'),
        ('sgp', 'Singapore', 'apac'),
        ('syd', 'Sydney', 'apac'),
        ('vin', 'Vint Hill', 'northAmerica'),
        ('waw', 'Warsaw', 'europe')
    )
    AS data(code, name, region)
    JOIN region ON region.code = data.region
    ORDER BY data.code
;

COMMIT;
