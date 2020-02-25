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

-- name: user-email
SELECT email
  FROM auth
  WHERE id = $1
;

-- name: user-exists
SELECT id
  FROM auth
  WHERE email = $1
;
