package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func hashPassword(password string) (string, error) {

	// cost parameters
	cost := 32768 // must be a power of two greater than 1
	// r and p must satisfy r * p < 2³⁰.
	r := 8
	p := 1

	hashKeyLenght := 32

	salt := make([]byte, hashKeyLenght)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, cost, r, p, hashKeyLenght)
	if err != nil {
		return "", err
	}

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func comparePasswords(storedPassword string, suppliedPassword string) (bool, error) {

	pwsalt := strings.Split(storedPassword, ".")

	if len(pwsalt) < 2 {
		return false, fmt.Errorf(storedPassword)
	}

	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
