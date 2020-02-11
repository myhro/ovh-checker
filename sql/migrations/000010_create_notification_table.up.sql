CREATE TABLE IF NOT EXISTS notification (
  id SERIAL PRIMARY KEY,
  auth_id INTEGER REFERENCES auth(id),
  server_id INTEGER REFERENCES server(id),
  country_id INTEGER REFERENCES country(id),
  sent_at TIMESTAMP WITH TIME ZONE,
  recurrent BOOLEAN DEFAULT false
);
