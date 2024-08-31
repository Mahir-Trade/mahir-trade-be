package queries

const (
	QueryCreateUser = `
		INSERT INTO users (email, phone_number, username, password, is_active) VALUES ($1, $2, $3, $4, TRUE) RETURNING user_id
	`

	QueryFindUserByEmailOrUsename = `
		SELECT user_id, uuid, phone_number, username, email, password FROM users WHERE email = $1 OR username = $1
	`

	QueryFindUserByEmailAndUsername = `
		SELECT user_id, uuid, phone_number, username, email, password FROM users WHERE email = $1 AND username = $2
	`

	QueryGetUserByID = `
		SELECT user_id, uuid, phone_number, username, email, password, is_active FROM users WHERE user_id = $1 AND deleted_at IS NULL
	`

	QueryGetUserByUUID = `
		SELECT user_id, uuid, phone_number, username, email, password FROM users WHERE uuid = $1 AND deleted_at IS NULL
	`

	QueryUpdateTypeUser = `
		UPDATE users
		SET is_active = $1, updated_by = $2, updated_at = NOW()
		WHERE user_id = $3
	`

	QueryGetUsers = `
		SELECT
			COUNT(*) OVER() as total_count,
			u.user_id,
			u.uuid,
			u.phone_number,
			u.email,
			u.username,
			u.is_active,
			case 
				when um.id is not null then 'Premium'
				else 'Standard'
			end as account_type,
			um.expired_at as membership_expired_date,
			u.created_at,
			u.created_by
		FROM
			users u
		left join
			user_memberships um on um.user_id = u.user_id AND um.expired_at > NOW()
		WHERE
			u.deleted_at IS NULL
	`

	QueryUpdatePassword = `
		UPDATE users
		SET password = $1, updated_by = $2, updated_at = NOW()
		WHERE user_id = $3
	`
)
