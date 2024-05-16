package queries

const (
	QueryCreatePackage = `
		INSERT INTO packages (price, duration_in_month, description, created_by, updated_by)
		VALUES ($1, $2, $3, 'SYSTEM', 'SYSTEM')
		RETURNING id
	`

	QueryGetPackages = `
		SELECT id, price, duration_in_month, description, created_by, updated_by, created_at, updated_at
		FROM packages
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	QueryGetTotalPackages = `
		SELECT COUNT(*) as total_count
		FROM packages
		WHERE deleted_at IS NULL
	`

	QueryGetPackageByID = `
		SELECT id, price, duration_in_month, description, created_by, updated_by, created_at, updated_at
		FROM packages
		WHERE id = $1
		AND deleted_at IS NULL
	`

	QueryUpdatePackage = `
		UPDATE packages
		SET price = $1, duration_in_month = $2, description = $3, updated_at = NOW(), updated_by = SYSTEM
		WHERE id = $4
	`

	QuertSoftDeletePackage = `
		UPDATE packages
		SET deleted_at = NOW(),
			deleted_by = $1,
			updated_at = NOW(),
			updated_by = $2
		WHERE id = $3
	`
)
