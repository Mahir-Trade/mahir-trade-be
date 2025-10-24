import { Request } from "express";

export interface FileType {
  type: string; // contoh: "png", "jpg"
  mimeType: string; // contoh: "image/png"
  size: number; // ukuran maksimal (MB)
}

/**
 * Validate file header (type, mime, size)
 */
export function validateFileHeader(
  file: Express.Multer.File,
  fileTypes: FileType[],
  keyName: string
): Error | null {
  const filenameSplit = file.originalname.split(".");
  const fileType = filenameSplit[filenameSplit.length - 1].toLowerCase();
  const mimeType = file.mimetype;

  for (const allowedType of fileTypes) {
    if (fileType === allowedType.type) {
      if (allowedType.size !== 0 && file.size > allowedType.size * (1 << 20)) {
        return new Error(
          `Error: Field validation for ${keyName}. File is too large`
        );
      }
      return null;
    }

    if (mimeType === allowedType.mimeType) {
      if (allowedType.size !== 0 && file.size > allowedType.size * (1 << 20)) {
        return new Error(
          `Error: Field validation for ${keyName}. File is too large`
        );
      }
      return null;
    }
  }

  return new Error(`Error: Field validation for ${keyName}. Invalid filetype`);
}

/**
 * Parse multipart form: return files & values
 * (gunakan middleware multer di route)
 */
export function parseMultipartForm(
  req: Request,
  fileTypes: FileType[]
): {
  files: { [fieldname: string]: Express.Multer.File[] };
  values: { [key: string]: string[] };
} {
  if (!req.files && !req.body) {
    throw new Error("Error: Field validation for parse multipart form");
  }

  const files: { [fieldname: string]: Express.Multer.File[] } = {};
  const values: { [key: string]: string[] } = {};

  // ambil semua file
  if (req.files) {
    if (Array.isArray(req.files)) {
      files["files"] = req.files;
    } else {
      for (const [key, val] of Object.entries(req.files)) {
        files[key] = val as Express.Multer.File[];
      }
    }
  }

  // ambil semua values (body text fields)
  for (const [key, val] of Object.entries(req.body)) {
    values[key] = Array.isArray(val) ? (val as string[]) : [val as string];
  }

  // validasi file types
  if (fileTypes.length > 0) {
    for (const [key, fileList] of Object.entries(files)) {
      for (const file of fileList) {
        const err = validateFileHeader(file, fileTypes, key);
        if (err) {
          throw new Error(`${key}, ${err.message}`);
        }
      }
    }
  }

  return { files, values };
}
