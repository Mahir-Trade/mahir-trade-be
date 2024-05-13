package queries

const (
	QueryCreateReport = `
		INSERT INTO reports (report_thumbnail_url, report_file_url, created_by, updated_by)
		VALUES ($1, $2, 'SYSTEM', 'SYSTEM')
		RETURNING id
	`

	QueryGetReports = `
		SELECT id, report_thumbnail_url, report_file_url, created_by, updated_by, created_at, updated_at
		FROM reports
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	QueryGetTotalReports = `
		SELECT COUNT(*) as total_count
		FROM reports
		WHERE deleted_at IS NULL
	`

	QueryGetReportByID = `
		SELECT id, report_thumbnail_url, report_file_url, created_by, updated_by, created_at, updated_at
		FROM reports
		WHERE id = $1
		AND deleted_at IS NULL
	`

	QueryUpdateReport = `
		UPDATE reports
		SET report_thumbnail_url = $1, report_file_url = $2, updated_at = NOW(), updated_by = 'SYSTEM'
		WHERE id = $3
	`

	QuertSoftDeleteReport = `
		UPDATE reports
		SET deleted_at = NOW(),
			deleted_by = $1,
			updated_at = NOW(),
			updated_by = $2
		WHERE id = $3
	`
)
