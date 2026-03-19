package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidatePassword_Valid(t *testing.T) {
	valid := []string{"Admin123!", "Passw0rd$", "Abc1234&", "Test55!x"}
	for _, pw := range valid {
		ok, err := ValidatePassword(pw)
		if !ok || err != nil {
			t.Errorf("ValidatePassword(%q) should be valid, got err: %v", pw, err)
		}
	}
}

func TestValidatePassword_TooShort(t *testing.T) {
	ok, err := ValidatePassword("Ab1!")
	if ok || err == nil {
		t.Error("expected short password to fail")
	}
}

func TestValidatePassword_NoLowercase(t *testing.T) {
	ok, err := ValidatePassword("ABCDEFG1!")
	if ok || err == nil {
		t.Error("expected no-lowercase to fail")
	}
}

func TestValidatePassword_NoUppercase(t *testing.T) {
	ok, err := ValidatePassword("abcdefg1!")
	if ok || err == nil {
		t.Error("expected no-uppercase to fail")
	}
}

func TestValidatePassword_NoDigit(t *testing.T) {
	ok, err := ValidatePassword("Abcdefgh!")
	if ok || err == nil {
		t.Error("expected no-digit to fail")
	}
}

func TestValidatePassword_NoSpecialChar(t *testing.T) {
	ok, err := ValidatePassword("Abcdefg1")
	if ok || err == nil {
		t.Error("expected no-special-char to fail")
	}
}

func TestHashPassword_ProducesValidBcrypt(t *testing.T) {
	hash, err := HashPassword("Test123!")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("Test123!")); err != nil {
		t.Error("bcrypt comparison should match")
	}
}

func TestHashPassword_DifferentHashesForSameInput(t *testing.T) {
	h1, _ := HashPassword("Test123!")
	h2, _ := HashPassword("Test123!")
	if h1 == h2 {
		t.Error("bcrypt should produce different hashes (different salts)")
	}
}

func TestHashPassword_WrongPasswordDoesNotMatch(t *testing.T) {
	hash, _ := HashPassword("Test123!")
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("Wrong123!")); err == nil {
		t.Error("wrong password should not match")
	}
}
