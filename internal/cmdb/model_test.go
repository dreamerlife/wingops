package cmdb

import "testing"

func TestIPFieldValidation(t *testing.T) {
	field := FieldDefinition{Name: "management_ip", Type: FieldTypeIP, Required: true}
	if err := field.Validate("10.0.1.10"); err != nil {
		t.Fatalf("expected valid ip: %v", err)
	}
	if err := field.Validate("not-an-ip"); err == nil {
		t.Fatal("expected invalid ip to fail")
	}
}
