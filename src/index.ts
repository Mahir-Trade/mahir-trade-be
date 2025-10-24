import cors from "cors";
import dotenv from "dotenv";
import express from "express";
import morgan from "morgan";
import expressListRoutes from "express-list-routes";
import router from "./routes/router";
import db from "./db/db.config";

dotenv.config();

const app = express();

// --- Middleware ---
app.use(express.json());
app.use(morgan("dev"));
app.use(
  cors({
    origin: [
      "http://localhost:3000",
      "http://localhost:3001",
      "http://localhost:3002",
      "http://localhost:3003",
    ],
    methods: ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization"],
    credentials: true,
  })
);

// --- Routes ---
app.use("/api/v1", router);

// --- Port dari ENV ---
const PORT = process.env.APP_PORT || 5000;

// --- Start server setelah DB terkoneksi ---
(async () => {
  try {
    const client = await db.connect();
    console.log("✅ Connected to PostgreSQL");
    client.release();

    app.listen(PORT, () => {
      console.log(`🚀 Server running on port ${PORT}`);
      expressListRoutes(router, { prefix: "/api/v1" });
    });
  } catch (err: any) {
    console.error("❌ Database connection failed:", err.message);
    process.exit(1);
  }
})();

export default app;
