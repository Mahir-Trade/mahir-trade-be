export const EmailQueries = {
  GetByKey: `
    SELECT body FROM email_templates
    WHERE key = $1
  `,
};
