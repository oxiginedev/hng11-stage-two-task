package models

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestSetPassword(t *testing.T) {
	user := &User{}
	password := "password"

	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Error setting password: %v", err)
	}

	if user.Password == "" {
		t.Fatal("Password hash should not be empty after setting password")
	}

	// Verify if the stored password hash matches the original password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		t.Fatalf("Password hash does not match the original password: %v", err)
	}
}

func TestComparePassword(t *testing.T) {
	user := &User{}
	password := "password"
	wrongPassword := "wrong-password"

	// Set password for the user
	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("Error setting password: %v", err)
	}

	// Test correct password
	match, err := user.ComparePassword(password)
	if err != nil {
		t.Fatalf("Error comparing password: %v", err)
	}
	if !match {
		t.Error("Expected password match, got mismatch")
	}

	// Test incorrect password
	match, err = user.ComparePassword(wrongPassword)
	if err != nil {
		t.Fatalf("Error comparing password: %v", err)
	}
	if match {
		t.Error("Expected password mismatch, got match")
	}
}
