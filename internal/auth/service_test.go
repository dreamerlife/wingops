package auth

import "testing"

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword(hash, "secret123") {
		t.Fatal("expected password to verify")
	}
	if VerifyPassword(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}
