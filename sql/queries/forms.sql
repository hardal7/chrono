-- name: CreateFeatureRequest :exec
INSERT INTO feature_requests (name, email, title, problem, feature, priority)
VALUES ($1, $2, $3, $4, $5, $6);
-- name: CreateBugReport :exec
INSERT INTO bug_reports (name, email, title, description, steps, environment, additional)
VALUES ($1, $2, $3, $4, $5, $6, $7);
