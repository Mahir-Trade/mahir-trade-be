package queries

const (
	QueryCreateSubModule                = `INSERT INTO sub_modules (module_id, sub_module_name, title, video_url, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	QueryCreateSubModuleWithoutModuleID = `INSERT INTO sub_modules (sub_module_name, title, video_url, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	QueryGetSubModules                  = `SELECT COUNT(id) OVER(), id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`

	QueryGetSubModuleByID = `SELECT id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE id = $1 AND deleted_at IS NULL`

	QueryGetSubModuleByModuleID = `SELECT id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE module_id = $1 AND deleted_at IS NULL`

	QueryUpdateSubModule = `UPDATE sub_modules SET sub_module_name = $1, title = $2, video_url = $3, updated_by = $4, updated_at = NOW(), module_id = $5 WHERE id = $6`

	QuerySoftDeleteSubModule = `UPDATE sub_modules SET deleted_at = NOW(), deleted_by = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3`

	QueryRemoveModuleIDFromSubModules = `
		UPDATE sub_modules
		SET
			module_id = NULL,
			updated_at = NOW(),
			updated_by = $1
		WHERE module_id = $2
	`
)
