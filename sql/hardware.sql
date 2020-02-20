-- name: latest-offers
SELECT
    hardware.id,
    server.name AS server,
    country.name AS country,
    MAX(offer.updated_at) AS updated_at
  FROM hardware
  JOIN server ON server.id = hardware.server_id
  JOIN offer ON hardware.id = offer.hardware_id
  JOIN datacenter ON datacenter.id = offer.datacenter_id
  JOIN country ON country.id = datacenter.country_id
  WHERE country.code = $1
    AND hardware.id BETWEEN $2 AND $3
  GROUP BY hardware.id, server.id, country.id
  ORDER BY hardware.id
