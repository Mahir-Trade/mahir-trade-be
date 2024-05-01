package queries

const (
	QueryCreateUser = `
		INSERT INTO users (email, fullname, phone_number, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id
	`
)