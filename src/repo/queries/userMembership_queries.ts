export const UserMembershipQueries = {
  CreateUserMembership: `
    INSERT INTO user_memberships 
      (user_id, package_id, exclusive_expired_at, is_membership_active, status, created_by, updated_by)
    VALUES 
      ($1, $2, $3, $4, $5, $6, $6)
    RETURNING id;
  `,

  UpdateUserMembershipExpired: `
    UPDATE user_memberships
    SET 
      expired_at = $1, 
      is_membership_active = $2, 
      status = $3, 
      updated_by = $4, 
      updated_at = NOW(), 
      exclusive_expired_at = $6, 
      package_id = $7
    WHERE user_id = $5
    RETURNING id;
  `,

  GetUserMembershipByUserID: `
    SELECT 
      id,
      user_id,
      package_id,
      COALESCE(expired_at, '1970-01-01') AS expired_at,
      COALESCE(exclusive_expired_at, '1970-01-01') AS exclusive_expired_at,
      is_membership_active,
      status,
      created_by,
      updated_by,
      created_at,
      updated_at
    FROM user_memberships
    WHERE user_id = $1
      AND deleted_at IS NULL;
  `,

  BulkUpdateUserMembershipExpired: `
    UPDATE user_memberships
    SET 
      is_membership_active = CASE
        WHEN expired_at < NOW() THEN FALSE
        ELSE is_membership_active
      END,
      updated_by = CASE
        WHEN expired_at < NOW() THEN 'CRONJOB'
        ELSE updated_by
      END;
  `,

  GetUserMemberships: `
    SELECT 
      id, 
      expired_at, 
      exclusive_expired_at 
    FROM user_memberships
    WHERE deleted_at IS NULL;
  `,

  GetUserMembershipExpired: `
    SELECT 
      id, 
      user_id, 
      expired_at, 
      is_membership_active, 
      COALESCE(exclusive_expired_at, '1970-01-01') AS exclusive_expired_at, 
      package_id
    FROM user_memberships
    WHERE (
      (exclusive_expired_at IS NOT NULL AND exclusive_expired_at < NOW())
      OR
      (exclusive_expired_at IS NULL AND expired_at < NOW())
    )
    AND is_membership_active IS TRUE
    AND deleted_at IS NULL;
  `,

  UpdateUserMembershipsByUserIDs: `
    UPDATE user_memberships
    SET 
      is_membership_active = FALSE, 
      status = 'EXPIRED', 
      updated_by = $1, 
      expired_at = NOW(), 
      exclusive_expired_at = NOW(), 
      updated_at = NOW();
  `,

  BulkUpdateMembershipPreOrderActivation: `
    UPDATE user_memberships um
    SET
      is_membership_active = TRUE,
      exclusive_expired_at = NOW() + (p.duration_in_month || ' months')::interval,
      updated_at = NOW(),
      updated_by = 'MEMBERSHIP_ACTIVATION'
    FROM packages p
    WHERE
      um.package_id = p.id
      AND um.status = 'PRE_ORDER'
      AND um.deleted_at IS NULL;
  `,
};
