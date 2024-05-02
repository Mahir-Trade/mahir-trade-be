package queries

const (
	QueryCreateUser = `
		INSERT INTO users (email, fullname, phone_number, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id
	`
	QueryFindUserByEmail = `
		SELECT user_id, uuid, fullname, phone_number, username, password FROM users WHERE email = $1 OR username = $1
	`
)
