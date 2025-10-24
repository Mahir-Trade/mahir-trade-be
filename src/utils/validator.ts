import { Request, Response, NextFunction } from "express";

export class Validator {
  static hasUppercase(value: string): boolean {
    return /[A-Z]/.test(value);
  }

  static hasLowercase(value: string): boolean {
    return /[a-z]/.test(value);
  }

  static hasNumber(value: string): boolean {
    return /[0-9]/.test(value);
  }

  static hasSpecialChar(value: string): boolean {
    return /[^a-zA-Z0-9]/.test(value);
  }
}

// Middleware generator mirip InitValidator
export function validatePasswordRules(field: string) {
  return (req: Request, res: Response, next: NextFunction) => {
    const value = req.body[field];

    if (typeof value !== "string") {
      return res.status(400).json({
        code: 400,
        message: `${field} must be a string`,
      });
    }

    if (!Validator.hasUppercase(value)) {
      return res.status(400).json({
        code: 400,
        message: `${field} must contain at least one uppercase letter`,
      });
    }

    if (!Validator.hasLowercase(value)) {
      return res.status(400).json({
        code: 400,
        message: `${field} must contain at least one lowercase letter`,
      });
    }

    if (!Validator.hasNumber(value)) {
      return res.status(400).json({
        code: 400,
        message: `${field} must contain at least one number`,
      });
    }

    if (!Validator.hasSpecialChar(value)) {
      return res.status(400).json({
        code: 400,
        message: `${field} must contain at least one special character`,
      });
    }

    next();
  };
}
