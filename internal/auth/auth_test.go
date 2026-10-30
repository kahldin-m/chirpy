package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


// TestHashPassword
func TestHashPassword(t *testing.T) {
	pass := "1234"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Errorf("Test failed: %v", err)
	}
	if hash == "" {
		t.Errorf("Expected non-empty hash")
	}
	t.Logf("hashing passed: %s\n", hash)
}


func TestPasswordHash(t *testing.T) {
	pass := "imahippo"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Errorf("Test failed: %v", err)
	}
	if hash == "" {
		t.Errorf("Expected non-empty hash")
	}

	match, err := CheckPasswordHash(pass, hash)
	if err != nil {
		t.Errorf("Test failed: %v", err)
	}
	if !match {
		t.Fatalf("expected password to match")
	}
}


func TestSamePasswordHash(t *testing.T) {
	hash1, _ := HashPassword("happyGOlucky")
	hash2, _ := HashPassword("happyGOlucky")

	t.Logf("hash1: %v\n", hash1)
	t.Logf("hash2: %v\n", hash2)
}



func TestCreateJWT(t *testing.T) {
	userID := uuid.New()
	jawsome, _ := MakeJWT(userID, "mushroom-sauce", time.Hour)

	t.Logf("userID: %v\ntoken: %v", userID, jawsome)
}



func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secretWord := "i-like-turtles"
	
	tokenString, err := MakeJWT(userID, secretWord, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotID, err := ValidateJWT(tokenString, secretWord) 
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotID != userID {
		t.Errorf("expected: %v\ngot: %v", userID, gotID)
	}
}
