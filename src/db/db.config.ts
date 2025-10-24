import { Pool } from "pg";
import dotenv from "dotenv";

dotenv.config();

function parseDuration(duration: string): number {
  if (!duration) return 0;
  const match = duration.match(/^(\d+)(s|m|h)?$/);
  if (!match) return parseInt(duration, 10);
  const value = parseInt(match[1], 10);
  const unit = match[2];
  switch (unit) {
    case "h":
      return value * 60 * 60 * 1000;
    case "m":
      return value * 60 * 1000;
    case "s":
      return value * 1000;
    default:
      return value;
  }
}

export interface DatabaseCfg {
  DBName: string;
  DBUser: string;
  DBPass: string;
  Host: string;
  Port: number;
  MaxOpenConns: number;
  MaxIdleConns: number;
  ConnMaxLifetime: number; // ms
}

// --- Load configuration dari ENV ---
const cfg: DatabaseCfg = {
  DBName: process.env.PG_DBNAME || process.env.DB_NAME || "mahir_trade",
  DBUser: process.env.PG_DBUSER || process.env.DB_USER || "MahirTrade",
  DBPass: process.env.PG_DBPASS || process.env.DB_PASS || "MahirTrade123-",
  Host: process.env.PG_HOST || process.env.DB_HOST || "mahir-trade-db",
  Port: parseInt(process.env.PG_PORT || process.env.DB_PORT || "5432", 10),

  MaxOpenConns: parseInt(process.env.PG_MAX_OPEN_CONNS || "30", 10),
  MaxIdleConns: parseInt(process.env.PG_MAX_IDLE_CONNS || "6", 10),
  ConnMaxLifetime: parseDuration(process.env.PG_CONN_MAX_LIFETIME || "30m"),
};

// --- Buat koneksi pool ke PostgreSQL ---
const db = new Pool({
  host: cfg.Host,
  port: cfg.Port,
  user: cfg.DBUser,
  password: cfg.DBPass,
  database: cfg.DBName,
  max: cfg.MaxOpenConns,
  idleTimeoutMillis: cfg.MaxIdleConns * 1000,
  connectionTimeoutMillis: 5000,
  ssl:
    process.env.PG_SSL_MODE === "require"
      ? { rejectUnauthorized: false }
      : false,
});

db.on("connect", () => {
  console.log("✅ Connected to PostgreSQL database:", cfg.DBName);
});

db.on("error", (err) => {
  console.error("💥 Unexpected PostgreSQL error:", err.message);
});

export default db;
