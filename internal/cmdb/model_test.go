package cmdb

import (
	"context"
	"testing"
)

func TestIPFieldValidation(t *testing.T) {
	field := FieldDefinition{Name: "management_ip", Type: FieldTypeIP, Required: true}
	if err := field.Validate("10.0.1.10"); err != nil {
		t.Fatalf("expected valid ip: %v", err)
	}
	if err := field.Validate("not-an-ip"); err == nil {
		t.Fatal("expected invalid ip to fail")
	}
}

func TestMemoryRepositoryStoresModelRelations(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	group, err := repo.CreateModelGroup(ctx, ModelGroup{Name: "middleware", DisplayName: "中间件"})
	if err != nil {
		t.Fatal(err)
	}
	server, err := repo.CreateModel(ctx, Model{GroupID: group.ID, Name: "server", DisplayName: "服务器"})
	if err != nil {
		t.Fatal(err)
	}
	mysql, err := repo.CreateModel(ctx, Model{GroupID: group.ID, Name: "mysql", DisplayName: "MySQL"})
	if err != nil {
		t.Fatal(err)
	}

	server.Relations = []ModelRelation{{
		TargetModelID: mysql.ID,
		RelationType:  "runs",
		DisplayName:   "运行",
	}}
	updated, err := repo.UpdateModel(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Relations) != 1 {
		t.Fatalf("expected one relation, got %#v", updated.Relations)
	}
	if updated.Relations[0].SourceModelID != server.ID || updated.Relations[0].TargetModelID != mysql.ID {
		t.Fatalf("relation endpoints did not round trip: %#v", updated.Relations[0])
	}
}
