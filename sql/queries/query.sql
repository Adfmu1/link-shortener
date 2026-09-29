-- name: InsertCode :one
INSERT INTO links (
    code, url
) VALUES (
    $1, $2
)
RETURNING code, created_at, click_count;

-- name: CheckIfCodeExistsFromUrl :one
SELECT code FROM links
WHERE url = $1;

-- name: GetDataFromCode :one
SELECT url, created_at FROM links
WHERE code = $1;

-- name: IncrementClicksFromCode :exec
UPDATE links
SET click_count = click_count + 1
WHERE code = $1;