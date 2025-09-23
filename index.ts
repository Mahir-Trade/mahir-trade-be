import cors from "cors";
import dotenv from "dotenv";
import express from "express";
import morgan from "morgan";
import db from "./db/db.config";
import router from "./routes/router";
import expressListRoutes from "express-list-routes";

dotenv.config();

const app = express();

// Middleware
app.use(express.json());
app.use(morgan("dev"));

app.use(
  cors({
    origin: ["http://localhost:3000"],
    credentials: true,
  })
);

// Routes
app.use("/api/v1", router);

const PORT = process.env.PORT || 5000;

// Database connection + start server
db.then(() => {
  app.listen(PORT, () => {
    console.log(`✅ Server is listening on port: ${PORT}`);

    // Print all routes clearly
    expressListRoutes(router, { prefix: "/api/v1" });
  });
});

export default app;
