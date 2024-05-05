package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

const saltSize = 16

// Generate 16 bytes randomely and securely using the
// Cryptographically secure pseudorandom number generator
func GenerateRandomSalt(saltSize int) ([]byte, error) {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])

	if err != nil {
		return nil, fmt.Errorf("failed generating new salt: %w", err)
	}

	return salt, nil
}

// Combine password and salt, then hash them using the SHA-512
// hashing algorithm and then return the hashed password as a
// hex string
func HashPassword(password string, salt []byte) string {
	// Convert password string to byte slice
	passwordBytes := []byte(password)

	// Create sha-512 hasher
	hasher := sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	hashedPasswordBytes := hasher.Sum(nil)

	// Convert the hashed password to a hex string
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func HashPasswordWithSalt(password string) (string, []byte, error) {
	salt, err := GenerateRandomSalt(24)
	if err != nil {
		return "", nil, fmt.Errorf("failed generating random salt when hashing password: %w", err)
	}

	hashedPassword := HashPassword(password, salt)

	return hashedPassword, salt, nil

}

// Check if the passwords match
func ComparePasswords(hashedPassword, currentPassword string, salt []byte) bool {
	currentPasswordHash := HashPassword(currentPassword, salt)

	return hashedPassword == currentPasswordHash
}
