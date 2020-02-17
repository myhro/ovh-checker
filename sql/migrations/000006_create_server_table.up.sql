BEGIN;

CREATE TABLE IF NOT EXISTS server (
  id SERIAL PRIMARY KEY,
  name TEXT,
  processor_id INTEGER REFERENCES processor(id),
  cores INTEGER,
  threads INTEGER,
  memory INTEGER,
  storage_id INTEGER REFERENCES storage(id),
  price DECIMAL
);

INSERT INTO server (name, processor_id, cores, threads, memory, storage_id, price)
  SELECT data.name, processor.id, data.cores, data.threads, data.memory, storage.id, data.price
    FROM (
      VALUES
        ('KS-1', 'Atom', 1, 2, 2, '500GB', 'HDD', 3.99),
        ('KS-2', 'Atom', 2, 4, 4, '1TB', 'HDD', 4.99),
        ('KS-3', 'Atom', 2, 4, 4, '2TB', 'HDD', 7.99),
        ('KS-4', 'Atom', 2, 4, 4, '2x2TB', 'HDD', 13.99),
        ('KS-5', 'Opteron', 4, 4, 16, '2TB', 'HDD', 13.99),
        ('KS-6', 'Core i5', 4, 4, 16, '2TB', 'HDD', 14.99),
        ('KS-7', 'Core i3', 2, 4, 8, '2TB', 'HDD', 14.99),
        ('KS-8', 'Core i7', 4, 8, 16, '2TB', 'HDD', 15.99),
        ('KS-9', 'Xeon', 4, 8, 16, '2x240GB', 'SSD', 16.99),
        ('KS-10', 'Core i5', 4, 4, 16, '2TB', 'HDD', 18.99),
        ('KS-11', 'Xeon', 4, 8, 16, '2x2TB', 'HDD', 19.99),
        ('KS-12', 'Xeon', 4, 8, 32, '2x2TB', 'HDD', 24.99)
    )
    AS data(name, processor, cores, threads, memory, disk_code, disk_type, price)
    JOIN processor ON processor.name = data.processor
    JOIN storage
      ON storage.code = data.disk_code
      AND storage.type = data.disk_type
    ORDER BY substring(data.name from '\d+')::INTEGER
;

COMMIT;
