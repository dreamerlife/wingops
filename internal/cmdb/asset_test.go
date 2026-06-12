package cmdb

import (
	"context"
	"testing"
)

func TestAssetValidatesRequiredFields(t *testing.T) {
	model := Model{
		Fields: []FieldDefinition{{Name: "hostname", Type: FieldTypeText, Required: true}},
	}
	asset := Asset{Attributes: map[string]any{}}
	if err := asset.Validate(model); err == nil {
		t.Fatal("expected missing hostname to fail")
	}
}

func TestMemoryRepositoryListsAssetsWithFiltersAndPagination(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	group, err := repo.CreateModelGroup(ctx, ModelGroup{Name: "compute", DisplayName: "计算资源"})
	if err != nil {
		t.Fatal(err)
	}
	model, err := repo.CreateModel(ctx, Model{
		GroupID:     group.ID,
		Name:        "server",
		DisplayName: "服务器",
		Fields: []FieldDefinition{
			{Name: "name", DisplayName: "名称", Type: FieldTypeText, Required: true},
			{Name: "management_ip", DisplayName: "管理 IP", Type: FieldTypeIP, Required: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	prod, err := repo.CreateAssetGroup(ctx, AssetGroup{Name: "prod", DisplayName: "生产环境", Dimension: "environment"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CreateAsset(ctx, Asset{
		ModelID:   model.ID,
		UniqueKey: "sn:web-01",
		Status:    AssetStatusRunning,
		GroupIDs:  []string{prod.ID},
		Attributes: map[string]any{
			"name":          "web-01",
			"management_ip": "10.0.0.11",
		},
	}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CreateAsset(ctx, Asset{
		ModelID:   model.ID,
		UniqueKey: "sn:db-01",
		Status:    AssetStatusMaintenance,
		Attributes: map[string]any{
			"name":          "db-01",
			"management_ip": "10.0.0.21",
		},
	}, "operator")
	if err != nil {
		t.Fatal(err)
	}

	result, err := repo.ListAssets(ctx, AssetListFilter{
		ModelID:  model.ID,
		GroupID:  prod.ID,
		Status:   AssetStatusRunning,
		Keyword:  "web",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || len(result.Items) != 1 {
		t.Fatalf("expected one filtered asset, got total=%d items=%d", result.Total, len(result.Items))
	}
	if result.Items[0].UniqueKey != "sn:web-01" {
		t.Fatalf("unexpected asset: %#v", result.Items[0])
	}
	if len(result.Items[0].GroupIDs) != 1 || result.Items[0].GroupIDs[0] != prod.ID {
		t.Fatalf("expected group membership to round trip: %#v", result.Items[0].GroupIDs)
	}
}

func TestMemoryRepositoryPreventsDeletingModelWithAssets(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	group, err := repo.CreateModelGroup(ctx, ModelGroup{Name: "compute", DisplayName: "计算资源"})
	if err != nil {
		t.Fatal(err)
	}
	model, err := repo.CreateModel(ctx, Model{
		GroupID:     group.ID,
		Name:        "server",
		DisplayName: "服务器",
		Fields: []FieldDefinition{
			{Name: "name", DisplayName: "名称", Type: FieldTypeText, Required: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.CreateAsset(ctx, Asset{
		ModelID:   model.ID,
		UniqueKey: "sn:web-01",
		Attributes: map[string]any{
			"name": "web-01",
		},
	}, "operator")
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.DeleteModel(ctx, model.ID); err != ErrModelHasAssets {
		t.Fatalf("expected ErrModelHasAssets, got %v", err)
	}
}
