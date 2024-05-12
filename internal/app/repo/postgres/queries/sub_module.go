package queries

const (
	QueryCreateSubModule = `INSERT INTO sub_modules (module_id, sub_module_name, title, video_url, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	QueryGetSubModules   = `SELECT id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`

	QueryGetSubModuleByID = `SELECT id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE id = $1 AND deleted_at IS NULL`

	QueryGetSubModuleByModuleID = `SELECT id, uuid, module_id, sub_module_name, title, video_url, created_by, updated_by, created_at, updated_at FROM sub_modules WHERE module_id = $1 AND deleted_at IS NULL`

	QueryUpdateSubModule = `UPDATE sub_modules SET sub_module_name = $1, title = $2, video_url = $3, updated_by = $4, updated_at = NOW() WHERE id = $5`

	QuerySoftDeleteSubModule = `UPDATE sub_modules SET deleted_at = NOW(), deleted_by = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3`
)
