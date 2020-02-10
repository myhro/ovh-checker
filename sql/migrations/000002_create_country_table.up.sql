BEGIN;

CREATE TABLE IF NOT EXISTS country (
  id SERIAL PRIMARY KEY,
  code TEXT,
  name TEXT,
  region_id INTEGER REFERENCES region(id)
);

INSERT INTO country (code, name, region_id)
  SELECT data.code, data.name, region.id
    FROM (
      VALUES
        ('au', 'Australia', 'apac'),
        ('ca', 'Canada', 'northAmerica'),
        ('de', 'Germany', 'europe'),
        ('fr', 'France', 'europe'),
        ('pl', 'Poland', 'europe'),
        ('sg', 'Singapore', 'apac'),
        ('uk', 'United Kingdom', 'europe'),
        ('us', 'United States', 'northAmerica')
    )
    AS data(code, name, region)
    JOIN region ON region.code = data.region
    ORDER BY data.code
;

COMMIT;
