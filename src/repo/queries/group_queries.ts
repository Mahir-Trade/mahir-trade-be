export const GroupQueries = {
  CreateGroup: `
    INSERT INTO groups (group_name, created_at, updated_at)
    VALUES ($1, NOW(), NOW())
    RETURNING id
  `,

  GetGroupByID: `
    SELECT id, uuid, group_name, created_at, updated_at
    FROM groups
    WHERE id = $1
    AND deleted_at IS NULL
  `,

  GetGroups: `
    SELECT id, uuid, group_name, created_at, updated_at
    FROM groups
    WHERE deleted_at IS NULL
  `,

  GetTotalGroups: `
    SELECT COUNT(*) as total_count
    FROM groups
    WHERE deleted_at IS NULL
  `,

  UpdateGroup: `
    UPDATE groups
    SET group_name = $1, updated_at = NOW()
    WHERE id = $2
  `,

  SoftDeleteGroup: `
    UPDATE groups
    SET deleted_at = NOW(),
        deleted_by = $1,
        updated_at = NOW(),
        updated_by = $2
    WHERE id = $3
  `,
};
