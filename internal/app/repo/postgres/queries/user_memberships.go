package queries

const (
	QueryCreateUserMembership = `
		INSERT INTO user_memberships (user_id, package_id, expired_at, is_membership_active, status, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id
	`

	QueryUpdateUserMembershipExpired = `
		UPDATE user_memberships
		SET expired_at = $1, is_membership_active = $2, status = $3, updated_by = $4, updated_at = NOW()
		WHERE user_id = $5
		RETURNING id
	`

	QueryGetUserMembershipByUserID = `
		SELECT id, user_id, package_id, expired_at, is_membership_active, status, created_by, updated_by, created_at, updated_at
		FROM user_memberships
		WHERE user_id = $1
			AND deleted_at IS NULL
	`
	QueryBulkUpdateUserMembershipExpired = `
		UPDATE user_memberships
		SET is_membership_active = CASE
			WHEN expired_at < NOW() THEN false
			ELSE is_membership_active
		END,
		updated_by = CASE
			WHEN expired_at < NOW() THEN 'CRONJOB'
			ELSE updated_by
		END
	`

	QueryGetUserMemberships = `
		SELECT id, expired_at FROM user_memberships
		WHERE deleted_at IS NULL
	`

	QueryGetUserMembershipExpired = `
		SELECT id, user_id, expired_at, is_membership_active FROM user_memberships
		WHERE expired_at < NOW()
			AND is_membership_active = true
			AND deleted_at IS NULL;
	`

	QueryUpdateUserMembershipsByUserIDs = `
		UPDATE user_memberships
		SET is_membership_active = false, status 'EXPIRED', updated_by = $1, expired_at = NOW(), updated_at = NOW()
	`

	QueryBulkUpdateMembershipPreOrderActivation = `
		UPDATE user_memberships um
		SET
			is_membership_active = TRUE,
			expired_at = now() + (p.duration_in_month || ' months')::interval,
			updated_at = now(),
			updated_by = 'MEMBERSHIP_ACTIVATION'
		FROM packages p
		WHERE
			um.package_id = p.id
			AND um.status = 'PRE_ORDER'
			AND um.deleted_at IS NULL;
		`
)
