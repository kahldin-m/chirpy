package auth

import (
	"testing"
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
	if match {
		t.Logf("check passed: %v", match)
	}
}


func TestSamePasswordHash(t *testing.T) {
	hash1, _ := HashPassword("happyGOlucky")
	hash2, _ := HashPassword("happyGOlucky")

	t.Logf("hash1: %v\n", hash1)
	t.Logf("hash2: %v\n", hash2)
}