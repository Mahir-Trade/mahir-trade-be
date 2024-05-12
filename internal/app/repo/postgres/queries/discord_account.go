package queries

const (
	QueryCreateDiscordAccount = `
		INSERT INTO discord_accounts (user_id, discord_account_id, username, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	QueryGetDiscordAccountByUserID = `
		SELECT id, user_id, discord_account_id, username, email, created_at, updated_at
		FROM discord_accounts
		WHERE user_id = $1
		AND deleted_at IS NULL
	`

	QuerySoftDeleteDiscordAccount = `
		UPDATE discord_accounts
		SET deleted_at = NOW(),
			deleted_by = 'SYSTEM',
			updated_at = NOW(),
			updated_by = 'SYSTEM'
		WHERE id = $1
	`
)
