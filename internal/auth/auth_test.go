package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "super_secret_password_123"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	// Check that the hash is not empty and not equal to the original
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == password {
		t.Fatal("hash should not equal the original password")
	}

	// Verify the correct password
	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error checking password: %v", err)
	}
	if !ok {
		t.Fatal("expected password to match hash")
	}

	// Verify an incorrect password
	ok, err = CheckPasswordHash("wrong_password", hash)
	if err != nil {
		t.Fatalf("unexpected error checking wrong password: %v", err)
	}
	if ok {
		t.Fatal("expected password check to fail for wrong password")
	}
}
