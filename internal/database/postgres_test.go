package database

import "testing"

func TestOpenRejectsEmptyDSN(t *testing.T) {
	if _, err := Open(""); err == nil {
		t.Fatal("expected error for empty dsn")
	}
}
