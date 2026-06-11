package cmdb

import "testing"

func TestAssetValidatesRequiredFields(t *testing.T) {
	model := Model{
		Fields: []FieldDefinition{{Name: "hostname", Type: FieldTypeText, Required: true}},
	}
	asset := Asset{Attributes: map[string]any{}}
	if err := asset.Validate(model); err == nil {
		t.Fatal("expected missing hostname to fail")
	}
}
