BEGIN;

CREATE TABLE IF NOT EXISTS datacenter (
  id SERIAL PRIMARY KEY,
  code TEXT,
  name TEXT,
  country_id INTEGER REFERENCES country(id)
);

INSERT INTO datacenter (code, name, country_id)
  SELECT data.code, data.name, country.id
    FROM (
      VALUES
        ('bhs', 'Beauharnois', 'ca'),
        ('default', 'Default', 'fr'),
        ('fra', 'Frankfurt', 'de'),
        ('gra', 'Gravelines', 'fr'),
        ('hil', 'Hillsboro', 'us'),
        ('lon', 'London', 'uk'),
        ('rbx', 'Roubaix', 'fr'),
        ('sbg', 'Strasbourg', 'fr'),
        ('sgp', 'Singapore', 'sg'),
        ('syd', 'Sydney', 'au'),
        ('vin', 'Vint Hill', 'us'),
        ('waw', 'Warsaw', 'pl')
    )
    AS data(code, name, country)
    JOIN country ON country.code = data.country
    ORDER BY data.code
;

COMMIT;
