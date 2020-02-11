-- name: add-notification
INSERT INTO notification (auth_id, server_id, country_id, recurrent)
  SELECT auth.id, server.id, country.id, recurrent
    FROM (
      VALUES
        ($1, $2, $3, $4::BOOLEAN)
    )
    AS data(email, server, country, recurrent)
    JOIN auth ON auth.email = data.email
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
    server.cores,
    server.threads,
    server.memory,
    FORMAT('%s %s', storage.code, storage.type) AS storage,
    country.name AS country,
    hardware.code AS hardware
  FROM notification
  JOIN auth ON auth.id = notification.auth_id
  JOIN datacenter ON datacenter.country_id = notification.country_id
  JOIN server ON server.id = notification.server_id
  JOIN storage ON storage.id = server.storage_id
  JOIN country ON country.id = notification.country_id
  JOIN hardware ON hardware.server_id = notification.server_id
  JOIN offer ON offer.hardware_id = hardware.id
  WHERE offer.status != 'unavailable'
    AND offer.datacenter_id = datacenter.id
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
