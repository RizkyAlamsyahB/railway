package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key-at-least-32-chars!!"
	userID := uuid.New()
	email := "test@example.com"
	role := "admin"

	token, err := GenerateToken(userID, email, role, nil, nil, secret, 24, "test-issuer")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("claims.Email = %v, want %v", claims.Email, email)
	}
	if claims.Role != "admin" {
		t.Errorf("claims.Role = %v, want %v", claims.Role, role)
	}
	if claims.VendorID != nil {
		t.Errorf("claims.VendorID = %v, want nil", claims.VendorID)
	}
	if claims.Subject != userID.String() {
		t.Errorf("claims.Subject = %v, want %v", claims.Subject, userID.String())
	}
	if claims.Issuer != "test-issuer" {
		t.Errorf("claims.Issuer = %v, want test-issuer", claims.Issuer)
	}
}

func TestGenerateAndValidateToken_WithVendorID(t *testing.T) {
	secret := "test-secret-key-at-least-32-chars!!"
	userID := uuid.New()
	vendorID := uuid.New()
	email := "vendor@example.com"
	role := "umkm"

	token, err := GenerateToken(userID, email, role, &vendorID, nil, secret, 24, "test-issuer")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.VendorID == nil {
		t.Fatal("claims.VendorID = nil, want non-nil")
	}
	if *claims.VendorID != vendorID {
		t.Errorf("claims.VendorID = %v, want %v", *claims.VendorID, vendorID)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, _ := GenerateToken(uuid.New(), "a@b.com", "customer", nil, nil, "secret-1", 24, "issuer")

	_, err := ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateToken() expected error for wrong secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	token, err := GenerateToken(uuid.New(), "a@b.com", "customer", nil, nil, "secret", -1, "issuer")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = ValidateToken(token, "secret")
	if err == nil {
		t.Fatal("ValidateToken() expected error for expired token")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	_, err := ValidateToken("not.a.valid.token", "secret")
	if err == nil {
		t.Fatal("ValidateToken() expected error for malformed token")
	}
}
