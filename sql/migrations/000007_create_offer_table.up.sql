CREATE TABLE IF NOT EXISTS offer (
  id BIGSERIAL PRIMARY KEY,
  hardware_id INTEGER REFERENCES hardware(id),
  datacenter_id INTEGER REFERENCES datacenter(id),
  status TEXT,
  updated_at TIMESTAMP WITH TIME ZONE,
  UNIQUE(hardware_id, datacenter_id)
);
