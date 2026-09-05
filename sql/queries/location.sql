-- name: SetLocationPrivacy :exec
UPDATE users
SET hide_country = $2, updated_at = NOW()
WHERE id = $1;
