import crypto from "crypto";

/**
 * AES-256-CBC Encryption with PKCS5 Padding
 */
export function encryptAES256CBC(
  plaintext: string,
  key: string,
  iv: string
): string {
  const bKey = Buffer.from(key, "utf8");

  console.log("bKey: ", bKey);

  const bIV = Buffer.from(iv, "utf8");
  console.log("bID: ", bIV);

  // PKCS5 padding (identik PKCS7 di Node.js, karena block size AES = 16)
  const cipher = crypto.createCipheriv("aes-256-cbc", bKey, bIV);
  let encrypted = cipher.update(plaintext, "utf8", "base64");
  encrypted += cipher.final("base64");

  return encrypted;
}

/**
 * AES-256-CBC Decryption with PKCS5 Unpadding
 */
export function decryptAES256CBC(
  ciphertext: string,
  key: string,
  iv: string
): string {
  const bKey = Buffer.from(key, "utf8");
  const bIV = Buffer.from(iv, "utf8");

  const decipher = crypto.createDecipheriv("aes-256-cbc", bKey, bIV);
  let decrypted = decipher.update(ciphertext, "base64", "utf8");
  decrypted += decipher.final("utf8");

  return decrypted;
}

/**
 * Manual PKCS5 Padding (opsional, kalau mau pakai buffer langsung)
 */
export function pkcs5Padding(data: Buffer, blockSize: number): Buffer {
  const padding = blockSize - (data.length % blockSize);
  const padtext = Buffer.alloc(padding, padding);
  return Buffer.concat([data, padtext]);
}

export function pkcs5Unpadding(data: Buffer): Buffer {
  if (data.length === 0) {
    throw new Error("data is empty");
  }
  const padding = data[data.length - 1];
  if (padding > data.length || padding === 0) {
    throw new Error("invalid padding");
  }
  return data.slice(0, data.length - padding);
}
