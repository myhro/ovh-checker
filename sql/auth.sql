-- name: add-user
INSERT INTO auth (email, password)
  VALUES ($1, crypt($2, gen_salt('bf')))
;

-- name: check-password
SELECT id
  FROM auth
  WHERE email = $1
    AND password = crypt($2, password)
;

-- name: user-exists
SELECT TRUE
  FROM auth
  WHERE email = $1
;
