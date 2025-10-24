import jwt, { SignOptions } from "jsonwebtoken";
import dotenv from "dotenv";
import { decryptAES256CBC, encryptAES256CBC } from "./security";
dotenv.config();

export interface Claims {
  data: any;
  iss?: string;
  exp?: number;
}

/**
 * Sama persis dengan Go: Sign
 */
export async function signJWT(
  data: any
): Promise<{ token: string; exp: Date }> {
  const secret = process.env.JWT_SECRET_KEY as string;
  if (!secret) throw new Error("JWT_SECRET_KEY is not set");

  // marshal data ke JSON string
  const payloadStr = JSON.stringify(data);

  console.log("JWT_ENCRYPT_KEY: ", process.env.JWT_ENCRYPT_KEY);
  console.log("JWT_ENCRYPT_IV: ", process.env.JWT_ENCRYPT_IV);

  // encrypt payload pakai AES
  const encryptedData = await encryptAES256CBC(
    payloadStr,
    process.env.JWT_ENCRYPT_KEY as string,
    process.env.JWT_ENCRYPT_IV as string
  );

  const exp = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000); // 7 days

  const claims: Claims = {
    data: encryptedData,
    iss: process.env.JWT_ISSUER,
    exp: Math.floor(exp.getTime() / 1000),
  };

  const token = jwt.sign(claims, secret, {
    algorithm: "HS256",
  } as SignOptions);

  return { token, exp };
}

/**
 * Sama persis dengan Go: Verify
 */
export async function verifyJWT(token: string): Promise<Claims> {
  const secret = process.env.JWT_SECRET_KEY as string;
  if (!secret) throw new Error("JWT_SECRET_KEY is not set");

  const decoded = jwt.verify(token, secret) as Claims;

  // decrypt data
  const decrypted = await decryptAES256CBC(
    decoded.data as string,
    process.env.JWT_ENCRYPT_KEY as string,
    process.env.JWT_ENCRYPT_IV as string
  );

  // parse kembali ke object
  decoded.data = JSON.parse(decrypted);

  return decoded;
}
