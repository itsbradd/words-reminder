# name: GetUser :one
SELECT * FROM user
WHERE id = ? LIMIT 1;

# name: SignUpUser :execlastid
INSERT INTO user (username, password)
VALUES (?, ?);