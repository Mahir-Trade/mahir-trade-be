import { body } from "express-validator";

// Register (kayak UserRegistrationRequest di Go)
export const registerValidator = [
  body("email").isEmail().withMessage("Email is not valid"),
  body("password")
    .isStrongPassword({
      minLength: 8,
      minLowercase: 1,
      minUppercase: 1,
      minNumbers: 1,
      minSymbols: 1,
    })
    .withMessage(
      "Password must contain uppercase, lowercase, number, and special char"
    ),
  body("username").notEmpty().withMessage("Username is required"),
];

// Login (kayak LoginReq di Go)
export const loginValidator = [
  body("identity")
    .notEmpty()
    .withMessage("Identity (email/username) is required"),
  body("password").notEmpty().withMessage("Password is required"),
];

// Forgot Password (kayak ForgotPasswordReq di Go)
export const forgotPasswordValidator = [
  body("email").isEmail().withMessage("Valid email is required"),
];

// Reset Password (kayak ResetPasswordRequest di Go)
export const resetPasswordValidator = [
  body("password")
    .isStrongPassword({
      minLength: 8,
      minLowercase: 1,
      minUppercase: 1,
      minNumbers: 1,
      minSymbols: 1,
    })
    .withMessage("Password does not meet security requirements"),
  body("password_confirmation")
    .custom((value, { req }) => value === req.body.password)
    .withMessage("Passwords do not match"),
];

// User Verification (kayak UserVerificationReq di Go)
export const userVerificationValidator = [
  body("email").notEmpty().withMessage("Email is required"),
  body("expired_time").notEmpty().withMessage("Expired time is required"),
];
