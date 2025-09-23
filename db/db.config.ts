import dotenv from "dotenv";
import mongoose from "mongoose";
dotenv.config();

const connection = `${process.env.DATABASE_URL}`;
const options = {
  dbName: process.env.DATABASE_NAME,
  autoIndex: false,
  maxPoolSize: 10,
  serverSelectionTimeoutMs: 5000,
  socketTimeoutMS: 45000,
  family: 4,
};

const db = mongoose
  .connect(connection, options)
  .then((ress) => {
    if (ress) {
      console.log(`Database connection successfully connected!`);
    }
  })
  .catch((err) => {
    console.log(err);
  });

export default db;
