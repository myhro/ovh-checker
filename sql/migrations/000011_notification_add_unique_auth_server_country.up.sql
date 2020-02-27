ALTER TABLE notification
  ADD CONSTRAINT notification_auth_id_server_id_country_id_key
  UNIQUE (auth_id, server_id, country_id)
;
