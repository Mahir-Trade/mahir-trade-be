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
		SELECT id, uuid, group_id, module_name, thumbnail_url, tag, created_by, created_at, updated_at, updated_by
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
		SELECT COUNT(*) OVER() as total_count, m.id, m.uuid, m.group_id, m.module_name, m.thumbnail_url, m.tag, m.created_by, m.created_at, m.updated_at, m.updated_by
		FROM modules AS m
		WHERE m.deleted_at IS NULL
		`

	QueryGetModulesWithSearch = `
		WITH count AS (
			SELECT COUNT(*) as total_count
			FROM modules
			WHERE module_name ILIKE '%' || $1 || '%' AND deleted_at IS NULL
		)
		SELECT count.total_count, m.id, m.uuid, m.group_id, m.module_name, m.thumbnail_url, m.tag, m.created_by, m.created_at, m.updated_at, m.updated_by
		FROM modules AS m
		JOIN count ON TRUE
		WHERE module_name ILIKE '%' || $1 || '%' AND deleted_at IS NULL
		ORDER BY m.created_at DESC
		LIMIT $2
		OFFSET $3
		`

	QueryGetModulesByGroupID = `
		SELECT m.id, m.uuid, m.group_id, m.module_name, m.thumbnail_url, m.created_by, m.created_at, m.updated_at, m.updated_by
		FROM modules AS m
		WHERE m.group_id = $1 AND m.deleted_at IS NULL
		`

	QuerySoftDeleteModule = `
		UPDATE modules
		SET deleted_at = NOW(), deleted_by = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3
		`

	QueryRemoveGroupIDFromModules = `
		UPDATE modules
		SET
			group_id = NULL,
			updated_at = NOW(),
			updated_by = $1
		WHERE group_id = $2
	`

	QueryGetPercentageModulesUser = `
	SELECT
    $1 AS user_id,
    $2 AS module_id,
    CASE
        WHEN COUNT(sm.id) = 0 THEN 0
        ELSE (COUNT(usm.id)::decimal / COUNT(sm.id)) * 100
    END AS completion_percentage
	FROM
		modules m
	LEFT JOIN
		sub_modules sm ON m.id = sm.module_id
	LEFT JOIN
		user_sub_modules usm ON sm.id = usm.sub_module_id AND usm.user_id = $3
	WHERE
		m.id = $4 AND sm.deleted_at IS NULL
	GROUP BY
		m.id
	ORDER BY
		m.id;
	 `
)
