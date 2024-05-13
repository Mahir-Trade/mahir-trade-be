package utils

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const (
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()-_=+,.?"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
func GenerateRandomPassword(length int) string {
	allChars := uppercaseLetters + lowercaseLetters + digits + specialChars
	password := make([]byte, length)
	password[0] = uppercaseLetters[rand.Intn(len(uppercaseLetters))]
	password[1] = lowercaseLetters[rand.Intn(len(lowercaseLetters))]
	password[2] = digits[rand.Intn(len(digits))]
	password[3] = specialChars[rand.Intn(len(specialChars))]
	for i := 4; i < length; i++ {
		password[i] = allChars[rand.Intn(len(allChars))]
	}
	for i := range password {
		j := rand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
