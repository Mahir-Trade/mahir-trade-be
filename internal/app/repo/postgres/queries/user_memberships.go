package queries

const (
	QueryCreateUserMembership = `
		INSERT INTO user_memberships (user_id, package_id, expired_at, is_membership_active, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`

	QueryUpdateUserMembershipExpired = `
		UPDATE user_memberships
		SET expired_at = $1, is_membership_active = $2, updated_by = $3
		WHERE user_id = $4
		RETURNING id
	`

	QueryGetUserMembershipByUserID = `
		SELECT id, user_id, package_id, expired_at, is_membership_active, created_by, updated_by, created_at, updated_at
		FROM user_memberships
		WHERE user_id = $1
			AND deleted_at IS NULL
	`
)
