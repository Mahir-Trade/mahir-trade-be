package queries

const (
	QueryCreateUser = `
		INSERT INTO users (email, fullname, phone_number, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id
	`
	QueryFindUserByEmailOrUsename = `
		SELECT user_id, uuid, fullname, phone_number, username, email, password FROM users WHERE email = $1 OR username = $1
	`

	QueryFindUserByEmailAndUsername = `
		SELECT user_id, uuid, fullname, phone_number, username, email, password FROM users WHERE email = $1 AND username = $2
	`

	QueryGetUserByID = `
		SELECT user_id, uuid, fullname, phone_number, username, email, password FROM users WHERE user_id = $1 AND deleted_at IS NULL
	`

	QueryGetUserByUUID = `
		SELECT user_id, uuid, fullname, phone_number, username, email, password FROM users WHERE uuid = $1 AND deleted_at IS NULL
	`

	QueryUpdateTypeUser = `
		UPDATE users
		SET is_active = $1, updated_by = $2, updated_at = NOW()
		WHERE user_id = $3
		`
)
