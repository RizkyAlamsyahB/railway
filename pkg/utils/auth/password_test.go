package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "securePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == password {
		t.Fatal("HashPassword() returned unhashed password")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("CheckPassword() with correct password error = %v", err)
	}

	if err := CheckPassword("wrongPassword", hash); err == nil {
		t.Error("CheckPassword() with wrong password expected error")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("HashPassword() produced identical hashes for same password (bcrypt should use random salt)")
	}
}
