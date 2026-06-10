package system

import "testing"

func TestConfigValueRequiresKey(t *testing.T) {
	config := Config{Key: "", Value: "enabled"}
	if config.Valid() {
		t.Fatal("expected config without key to be invalid")
	}
}
