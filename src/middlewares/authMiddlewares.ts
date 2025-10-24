import { Request, Response, NextFunction } from "express";
import jwt from "jsonwebtoken";
import crypto from "crypto";
import { UserContext } from "../config/userContext";

export const authMiddleware = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  try {
    const authHeader = req.headers.authorization;

    if (!authHeader || !authHeader.startsWith("Bearer ")) {
      console.error(
        "❌ [AuthMiddleware] Header Authorization kosong / format salah"
      );
      return res.status(401).json({
        code: 401,
        message: "Unauthorized - Missing Bearer token",
      });
    }

    const token = authHeader.split(" ")[1];
    const secretKey = process.env.JWT_SECRET_KEY || "secret";

    const decoded: any = jwt.verify(token, secretKey);
    console.log("✅ [AuthMiddleware] Token ter-decode:", decoded);

    // 🔓 Decrypt payload
    const key = Buffer.from(process.env.JWT_ENCRYPT_KEY!, "utf-8");
    const iv = Buffer.from(process.env.JWT_ENCRYPT_IV!, "utf-8");

    const decipher = crypto.createDecipheriv("aes-256-cbc", key, iv);
    let decrypted = decipher.update(decoded.data, "base64", "utf8");
    decrypted += decipher.final("utf8");

    const payload = JSON.parse(decrypted);
    console.log("🧩 [AuthMiddleware] Payload setelah dekripsi:", payload);

    // 🔥 Simpan ke dua tempat sekaligus
    (req as any).userData = payload; // versi lama tetap ada
    UserContext.run(payload, () => {
      console.log("🔑 [DEBUG] Context injected untuk:", payload.username);
      next();
    });
  } catch (err: any) {
    console.error("🔥 [AuthMiddleware] Error:", err.message);
    return res.status(401).json({
      code: 401,
      message: "Invalid or expired token",
    });
  }
};
