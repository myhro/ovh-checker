BEGIN;

CREATE TABLE IF NOT EXISTS hardware (
  id SERIAL PRIMARY KEY,
  code TEXT,
  server_id INTEGER REFERENCES server(id)
);

INSERT INTO hardware (code, server_id)
  SELECT data.code, server.id
    FROM (
      VALUES
        -- Europe
        ('1801sk12', 'KS-1'),
        ('1801sk13', 'KS-2'),
        ('1801sk14', 'KS-3'),
        ('1801sk15', 'KS-4'),
        ('1801sk16', 'KS-5'),
        ('1801sk17', 'KS-6'),
        ('1801sk18', 'KS-7'),
        ('1801sk19', 'KS-8'),
        ('1801sk20', 'KS-9'),
        ('1801sk21', 'KS-10'),
        ('1801sk22', 'KS-11'),
        ('1801sk23', 'KS-12'),
        -- North America
        ('1804sk12', 'KS-1'),
        ('1804sk16', 'KS-5'),
        ('1804sk18', 'KS-7'),
        ('1804sk20', 'KS-9'),
        ('1804sk21', 'KS-10'),
        ('1804sk22', 'KS-11'),
        ('1804sk23', 'KS-12')
    )
    AS data(code, server_name)
    JOIN server ON server.name = data.server_name
    ORDER BY data.code
;

COMMIT;
