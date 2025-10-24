import { Request, Response, NextFunction } from "express";
import { validationResult, ValidationError } from "express-validator";

export const validateRequest = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({
      code: 400,
      message: "Invalid request body",
      errors: errors.array().map((err: ValidationError) => {
        if (err.type === "field") {
          return { field: err.path, msg: err.msg };
        }
        return { field: "unknown", msg: err.msg };
      }),
    });
  }
  next();
};
