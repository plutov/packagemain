-- name: GetUsersWithMinPosts :many
SELECT u.name, COUNT(p.id) AS post_count
FROM users AS u
JOIN posts AS p ON u.id = p.user_id
GROUP BY u.id
HAVING COUNT(p.id) >= 2;
