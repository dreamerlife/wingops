package auth

import (
	"context"
	"testing"
)

func TestRoleAllowsPermission(t *testing.T) {
	role := Role{Name: "readonly", Permissions: []Permission{{Code: "cmdb.asset.read"}}}
	if !role.Allows("cmdb.asset.read") {
		t.Fatal("expected read permission")
	}
	if role.Allows("cmdb.asset.write") {
		t.Fatal("did not expect write permission")
	}
}

func TestMemoryRepositoryManagesUsersAndRoles(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	if _, err := repo.CreateRole(ctx, Role{
		Name:        "ops_viewer",
		DisplayName: "运维观察员",
		Permissions: []Permission{
			{Code: "cmdb.asset.read"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	user, err := repo.CreateUser(ctx, User{
		Username:    "liwei",
		DisplayName: "李伟",
		Status:      "active",
	}, "secret123", []string{"ops_viewer"})
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == "" || len(user.Roles) != 1 || user.Roles[0].Name != "ops_viewer" {
		t.Fatalf("unexpected user: %#v", user)
	}

	updated, err := repo.UpdateRole(ctx, Role{
		Name:        "ops_viewer",
		DisplayName: "运维只读",
		Permissions: []Permission{
			{Code: "cmdb.asset.read"},
			{Code: "cmdb.model.read"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Allows("cmdb.model.read") {
		t.Fatalf("expected updated permission: %#v", updated)
	}

	user.Status = "disabled"
	user.DisplayName = "李伟-停用"
	user.Roles = []Role{{Name: "ops_viewer"}}
	saved, err := repo.UpdateUser(ctx, user, "")
	if err != nil {
		t.Fatal(err)
	}
	if saved.Active() || saved.DisplayName != "李伟-停用" {
		t.Fatalf("unexpected updated user: %#v", saved)
	}
}
