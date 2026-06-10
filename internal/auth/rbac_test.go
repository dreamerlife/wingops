package auth

import "testing"

func TestRoleAllowsPermission(t *testing.T) {
	role := Role{Name: "readonly", Permissions: []Permission{{Code: "cmdb.asset.read"}}}
	if !role.Allows("cmdb.asset.read") {
		t.Fatal("expected read permission")
	}
	if role.Allows("cmdb.asset.write") {
		t.Fatal("did not expect write permission")
	}
}
