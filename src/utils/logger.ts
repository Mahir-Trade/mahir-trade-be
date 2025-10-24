import fs from "fs";
import path from "path";

/**
 * Simple structured logger
 * Inspired by Go slog.Logger
 *
 * Usage:
 *  Logger.info("Server started", { port: 8080 })
 *  Logger.error("DB connection failed", err)
 */

export class Logger {
  private static logDir = path.resolve("logs");
  private static logFile = path.join(Logger.logDir, "app.log");

  private static ensureLogDir() {
    if (!fs.existsSync(Logger.logDir)) {
      fs.mkdirSync(Logger.logDir, { recursive: true });
    }
  }

  private static write(level: string, message: string, meta?: any) {
    Logger.ensureLogDir();

    const timestamp = new Date().toISOString();
    const formatted =
      `[${timestamp}] [${level.toUpperCase()}] ${message}` +
      (meta ? ` | ${JSON.stringify(meta, null, 2)}` : "");

    // print to console
    if (level === "error") {
      console.error(formatted);
    } else if (level === "warn") {
      console.warn(formatted);
    } else {
      console.log(formatted);
    }

    // append to log file
    fs.appendFile(Logger.logFile, formatted + "\n", (err) => {
      if (err) console.error("Failed to write log:", err);
    });
  }

  static info(message: string, meta?: any) {
    Logger.write("info", message, meta);
  }

  static warn(message: string, meta?: any) {
    Logger.write("warn", message, meta);
  }

  static error(message: string, meta?: any) {
    Logger.write("error", message, meta);
  }

  static debug(message: string, meta?: any) {
    if (process.env.NODE_ENV !== "production") {
      Logger.write("debug", message, meta);
    }
  }
}
