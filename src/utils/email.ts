import Handlebars from "handlebars";

/**
 * Sama persis kayak MappingValuesToTemplate di Go
 * Parse template string lalu inject values
 */
export function mappingValuesToTemplate(
  values: Record<string, any>,
  template: string
): string {
  try {
    const compiled = Handlebars.compile(template);
    return compiled(values);
  } catch (err: any) {
    throw new Error(`failed to parse template: ${err.message}`);
  }
}
