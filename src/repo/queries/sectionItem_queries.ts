export const SectionItemQueries = {
  GetAllSectionItems: `
    SELECT id, section_id, title, subtitle, subjek, image_url, icon_url, "order", extra_data, created_at, updated_at
    FROM section_items
    ORDER BY section_id, "order" ASC;
  `,

  InsertSectionItem: `
    INSERT INTO section_items (
      section_id, title, subtitle, subjek, image_url, icon_url, "order", extra_data
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb)
    RETURNING id, section_id, title, subtitle, subjek, image_url, icon_url, "order", extra_data, created_at, updated_at;
  `,

  GetItemsBySection: `
    SELECT id, section_id, title, subtitle, image_url, icon_url, "order", extra_data, created_at, updated_at
    FROM section_items
    WHERE section_id = $1
    ORDER BY "order" ASC;
  `,

  UpdateSectionItem: `
    UPDATE section_items
    SET title = $1, subtitle = $2, subjek = $3, image_url = $4, icon_url = $5, "order" = $6, extra_data = $7::jsonb, updated_at = now()
    WHERE id = $8;
  `,

  DeleteSectionItem: `
    DELETE FROM section_items WHERE id = $1;
  `,

  GetSectionItemByID: `
    SELECT id, section_id, title, subtitle, image_url, icon_url, "order", extra_data, created_at, updated_at
    FROM section_items
    WHERE id = $1;
  `,

  CheckDuplicateSectionItemTitle: `
    SELECT EXISTS (
      SELECT 1 FROM section_items
      WHERE section_id = $1 AND title = $2
    );
  `,
};
