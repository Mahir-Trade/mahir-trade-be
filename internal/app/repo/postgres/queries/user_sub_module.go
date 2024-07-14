package queries

const (
	QueryCreateUserSubModule = `INSERT INTO user_sub_modules (user_id, sub_module_id, created_by) VALUES ($1, $2, $3) RETURNING id`

	QueryGetUserSubModuleBySubModuleIDAndUserID = `SELECT id, uuid, user_id, sub_module_id, created_by, updated_by, created_at, updated_at FROM user_sub_modules WHERE sub_module_id = $1 AND user_id = $2 AND deleted_at IS NULL`
)
