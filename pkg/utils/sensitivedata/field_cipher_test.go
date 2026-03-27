package sensitivedata

import "testing"

func TestFieldCipherRoundTrip(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	cipher, err := NewFieldCipher(key)
	if err != nil {
		t.Fatalf("NewFieldCipher() unexpected error: %v", err)
	}

	encrypted, err := cipher.EncryptString("1234567890123456")
	if err != nil {
		t.Fatalf("EncryptString() unexpected error: %v", err)
	}
	if encrypted == "1234567890123456" {
		t.Fatal("expected encrypted value to differ from plaintext")
	}

	decrypted, err := cipher.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString() unexpected error: %v", err)
	}
	if decrypted != "1234567890123456" {
		t.Fatalf("DecryptString() = %q, want %q", decrypted, "1234567890123456")
	}
}

func TestFieldCipherDecryptLegacyPlaintext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	cipher, err := NewFieldCipher(key)
	if err != nil {
		t.Fatalf("NewFieldCipher() unexpected error: %v", err)
	}

	decrypted, err := cipher.DecryptString("legacy-plaintext")
	if err != nil {
		t.Fatalf("DecryptString() unexpected error: %v", err)
	}
	if decrypted != "legacy-plaintext" {
		t.Fatalf("DecryptString() = %q, want legacy plaintext passthrough", decrypted)
	}
}
