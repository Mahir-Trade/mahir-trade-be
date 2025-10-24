import bcrypt from "bcryptjs";

/**
 * Constants sama dengan Go
 */
const uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
const lowercaseLetters = "abcdefghijklmnopqrstuvwxyz";
const digits = "0123456789";
const specialChars = "!@#$%^&*()-_=+,.?";

/**
 * Hash password pakai bcrypt (DefaultCost = 10)
 */
export async function hashPassword(password: string): Promise<string> {
  const salt = await bcrypt.genSalt(10);
  const hashed = await bcrypt.hash(password, salt);
  return hashed;
}

/**
 * Verifikasi password dengan hashed
 */
export async function verifyPassword(
  password: string,
  hashedPassword: string
): Promise<boolean> {
  return bcrypt.compare(password, hashedPassword);
}

/**
 * Generate random password dengan kombinasi
 * uppercase, lowercase, digit, special char
 */
export function generateRandomPassword(length: number): string {
  const allChars = uppercaseLetters + lowercaseLetters + digits + specialChars;

  const password: string[] = new Array(length);

  // Pastikan minimal ada 1 char dari tiap kategori
  password[0] =
    uppercaseLetters[Math.floor(Math.random() * uppercaseLetters.length)];
  password[1] =
    lowercaseLetters[Math.floor(Math.random() * lowercaseLetters.length)];
  password[2] = digits[Math.floor(Math.random() * digits.length)];
  password[3] = specialChars[Math.floor(Math.random() * specialChars.length)];

  for (let i = 4; i < length; i++) {
    password[i] = allChars[Math.floor(Math.random() * allChars.length)];
  }

  // Shuffle biar random
  for (let i = password.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [password[i], password[j]] = [password[j], password[i]];
  }

  return password.join("");
}
