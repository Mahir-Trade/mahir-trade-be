export const PackageQueries = {
  CreatePackage: `
    INSERT INTO packages (price, duration_in_month, description, discounted_price, discount_expired, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $6)
    RETURNING id
  `,

  GetPackages: `
    SELECT COUNT(*) OVER() as total_count, id, price, duration_in_month, description, discounted_price, discount_expired, created_by, updated_by, created_at, updated_at
    FROM packages
    WHERE deleted_at IS NULL
    LIMIT $1 OFFSET $2
  `,

  GetPackageByID: `
    SELECT id, price, duration_in_month, description, discounted_price, discount_expired, created_by, updated_by, created_at, updated_at
    FROM packages
    WHERE id = $1
    AND deleted_at IS NULL
  `,

  UpdatePackage: `
    UPDATE packages
    SET price = $1, duration_in_month = $2, description = $3, discounted_price = $4, updated_at = NOW(), updated_by = $6
    WHERE id = $5
  `,

  SoftDeletePackage: `
    UPDATE packages
    SET deleted_at = NOW(),
        deleted_by = $1,
        updated_at = NOW(),
        updated_by = $2
    WHERE id = $3
  `,

  UpdatePackageDiscountExpired: `
    UPDATE packages
    SET discount_expired = $1
  `,
};
