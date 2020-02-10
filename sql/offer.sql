-- name: available
SELECT DISTINCT ON (hardware.code) offer.id, country.name AS country, server.name AS server, hardware.code
  FROM offer
  JOIN hardware ON hardware.id = hardware_id
  JOIN datacenter ON datacenter.id = datacenter_id
  JOIN server ON server.id = hardware.server_id
  JOIN country ON country.id = datacenter.country_id
  WHERE offer.status != 'unavailable'
  ORDER BY hardware.code
;

-- name: import-json
WITH
  raw_input(body) AS (
    VALUES($1::json)
  ),
  expand_list AS (
    SELECT json_array_elements(body) AS elem
      FROM raw_input
  ),
  extract_fields AS (
    SELECT elem->>'hardware' AS hardware, elem->>'region' AS region, elem->>'datacenters' AS datacenters
      FROM expand_list
  ),
  expand_datacenters AS (
    SELECT hardware, region, json_array_elements(datacenters::json) AS dc
      FROM extract_fields
  )

INSERT INTO offer (hardware_id, datacenter_id, status)
  SELECT hardware.id, datacenter.id, data.availability
    FROM (
      SELECT hardware, dc->>'datacenter' AS datacenter, dc->>'availability' AS availability
        FROM expand_datacenters
    )
    AS data(hardware, datacenter, availability)
    JOIN hardware ON hardware.code = data.hardware
    JOIN datacenter ON datacenter.code = data.datacenter
    EXCEPT
      SELECT hardware_id, datacenter_id, status
        FROM offer
  ON CONFLICT (datacenter_id, hardware_id)
    DO UPDATE
      SET status = EXCLUDED.status, updated_at = NOW()
      WHERE offer.status != EXCLUDED.status
;
