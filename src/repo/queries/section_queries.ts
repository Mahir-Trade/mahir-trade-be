export const SectionQueries = {
  InsertSection: `
    INSERT INTO sections (slug, type, title, subtitle, "order")
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (slug) DO UPDATE
    SET 
      type = EXCLUDED.type,
      title = EXCLUDED.title,
      subtitle = EXCLUDED.subtitle,
      "order" = EXCLUDED."order",
      updated_at = now();
  `,

  GetSectionByType: `
    SELECT id, slug, type, title, subtitle, "order", created_at, updated_at
    FROM sections
    WHERE type = $1
    LIMIT 1;
  `,

  GetSectionBySlug: `
    SELECT id, slug, type, title, subtitle, "order", created_at, updated_at
    FROM sections
    WHERE slug = $1;
  `,

  UpdateSection: `
    UPDATE sections
    SET slug = $1,
        type = $2,
        title = $3,
        subtitle = $4,
        "order" = $5,
        updated_at = now()
    WHERE id = $6;
  `,

  GetSectionsWithItems: `
    SELECT 
      s.id, s.slug, s.type, s.title, s.subtitle, s."order",
      i.id, i.section_id, i.title, i.subtitle, i.image_url, i.icon_url, i."order", i.extra_data
    FROM sections s
    LEFT JOIN section_items i ON s.id = i.section_id
    ORDER BY s."order" ASC, i."order" ASC;
  `,

  FindSectionByID: `
    SELECT id, title, subtitle, created_at, updated_at
    FROM sections
    WHERE id = $1
    LIMIT 1;
  `,

  GetAllSectionsWithItems: `
  SELECT 
    s.id AS section_id,
    s.slug AS section_slug,
    s.type AS section_type,
    s.title AS section_title,
    s.subtitle AS section_subtitle,
    s."order" AS section_order,
    s.created_at AS section_created_at,
    s.updated_at AS section_updated_at,

    si.id AS item_id,
    si.title AS item_title,
    si.subtitle AS item_subtitle,
    si.subjek AS item_subjek,
    si.image_url,
    si.icon_url,
    si."order" AS item_order,
    si.extra_data,
    si.created_at AS item_created_at,
    si.updated_at AS item_updated_at

  FROM sections s
  LEFT JOIN section_items si ON si.section_id = s.id
  ORDER BY s."order" ASC, si."order" ASC;
`,

  GetSectionByID: `
    SELECT id, slug, type, title, subtitle, "order", created_at, updated_at
    FROM sections
    WHERE id = $1;
  `,

  GetSectionWithItemsByID: `
    SELECT
      s.id,
      s.slug,
      s.type,
      s.title,
      s.subtitle,
      s."order",
      s.created_at,
      s.updated_at,
      si.id,
      si.section_id,
      si.title,
      si.subtitle,
      si.image_url,
      si.icon_url,
      si."order",
      si.extra_data,
      si.created_at,
      si.updated_at
    FROM sections s
    LEFT JOIN section_items si ON s.id = si.section_id
    WHERE s.id = $1
    ORDER BY si."order" ASC;
  `,

  DeleteSection: `
    DELETE FROM sections WHERE id = $1;
  `,
};
