-- name: add-notification
INSERT INTO notification (auth_id, server_id, country_id, recurrent)
  SELECT auth_id, server.id, country.id, recurrent
    FROM (
      VALUES
        ($1::INTEGER, $2, $3, $4::BOOLEAN)
    )
    AS data(auth_id, server, country, recurrent)
    JOIN server ON server.name = data.server
    JOIN country ON country.code = data.country
;

-- name: mark-as-sent
UPDATE notification
  SET sent_at = $1
  WHERE id = $2
;

-- name: pending-notifications
SELECT DISTINCT ON (notification.id)
    notification.id,
    auth.email,
    server.name AS server,
    processor.name AS processor,
    server.cores,
    server.threads,
    server.memory,
    FORMAT('%s %s', storage.code, storage.type) AS storage,
    country.name AS country,
    hardware.code AS hardware
  FROM notification
  JOIN auth ON auth.id = notification.auth_id
  JOIN country ON country.id = notification.country_id
  JOIN datacenter ON datacenter.country_id = country.id
  JOIN server ON server.id = notification.server_id
  JOIN processor ON processor.id = server.processor_id
  JOIN storage ON storage.id = server.storage_id
  JOIN offer ON offer.datacenter_id = datacenter.id
  JOIN hardware ON hardware.id = offer.hardware_id
  WHERE offer.status != 'unavailable'
    AND notification.server_id = hardware.server_id
    AND (
      notification.sent_at IS NULL
      OR (
        notification.recurrent = true
        AND
        notification.sent_at < COALESCE(offer.updated_at, to_timestamp(0))
      )
    )
  ORDER BY notification.id
;
