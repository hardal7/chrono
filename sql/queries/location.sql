-- name: SetLocationPrivacy :exec
UPDATE users
SET hide_country = $2, updated_at = now()
WHERE id = $1;
