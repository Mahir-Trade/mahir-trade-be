package queries

const (
	QueryCreateModuleWithGroupID = `
		INSERT INTO modules (group_id, module_name, thumbnail_url, created_by) VALUES ($1, $2, $3, $4) RETURNING id
		`
	QueryCreateModuleWithGroupIDAndTag = `
		INSERT INTO modules (group_id, module_name, thumbnail_url, tag, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id
		`
	QueryCreateModuleWithoutGroupID = `
		INSERT INTO modules (module_name, thumbnail_url, tag, created_by) VALUES ($1, $2, $3, $4) RETURNING id
		`
	QueryCreateModuleWithoutGroupIDAndTag = `
		INSERT INTO modules (module_name, thumbnail_url, created_by) VALUES ($1, $2, $3) RETURNING id
		`

	QueryGetModuleByID = `
		SELECT id, uuid, group_id, module_name, thumbnail_url, tag, created_by, created_at, updated_at
		FROM modules
		WHERE id = $1
		`

	QueryUpdateModule = `
		UPDATE modules
		SET module_name = $1, updated_by = $2, updated_at = NOW()
		WHERE id = $3
		`

	QueryUpdateModuleWithThumbnail = `
		UPDATE modules
		SET module_name = $1, thumbnail_url = $2, updated_by = $3, updated_at = NOW()
		WHERE id = $4
		`

	QueryUpdateModuleWithTag = `
		UPDATE modules
		SET module_name = $1, tag = $2, updated_by = $3, updated_at = NOW()
		WHERE id = $4
		`

	QueryUpdateModuleWithThumbnailAndTag = `
		UPDATE modules
		SET module_name = $1, thumbnail_url = $2, tag = $3, updated_by = $4, updated_at = NOW()
		WHERE id = $5
		`

	QueryGetModules = `
		WITH count AS (
		SELECT COUNT(*) as total_count
		FROM modules
		)
		SELECT count.total_count, m.id, m.group_id, m.module_name, m.thumbnail_url, m.created_by, m.created_at, m.updated_at
		FROM modules AS m
		JOIN count ON TRUE
		LIMIT $1
		OFFSET $2
		`

	QueryGetModulesWithSearch = `
		WITH count AS (
			SELECT COUNT(*) as total_count
			FROM modules
			WHERE module_name ILIKE '%' || $1 || '%'
		)
		SELECT count.total_count, modules.id, modules.group_id, modules.module_name, modules.thumbnail_url, modules.created_by, modules.created_at, modules.updated_at
		FROM modules
		JOIN count ON TRUE
		WHERE module_name ILIKE '%' || $1 || '%'
		LIMIT $2
		OFFSET $3
		`
)
