package cmdb

import (
	"errors"
	"time"
)

const (
	AssetStatusPurchased   = "purchased"
	AssetStatusRacked      = "racked"
	AssetStatusRunning     = "running"
	AssetStatusMaintenance = "maintenance"
	AssetStatusRetired     = "retired"
)

type Asset struct {
	ID         string         `json:"id"`
	ModelID    string         `json:"model_id"`
	UniqueKey  string         `json:"unique_key"`
	Status     string         `json:"status"`
	GroupIDs   []string       `json:"group_ids,omitempty"`
	Attributes map[string]any `json:"attributes"`
}

type AssetGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Dimension   string `json:"dimension"`
}

type AssetListFilter struct {
	ModelID  string
	GroupID  string
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type AssetListResult struct {
	Items    []Asset `json:"items"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

type AssetChangeLog struct {
	ID          string         `json:"id"`
	AssetID     string         `json:"asset_id"`
	ActorID     string         `json:"actor_id"`
	BeforeValue map[string]any `json:"before_value"`
	AfterValue  map[string]any `json:"after_value"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (a Asset) Validate(model Model) error {
	if a.UniqueKey == "" {
		return errors.New("asset unique_key is required")
	}
	if a.Status != "" && !validAssetStatus(a.Status) {
		return errors.New("invalid asset status")
	}
	for _, field := range model.Fields {
		if err := field.Validate(a.Attributes[field.Name]); err != nil {
			return err
		}
	}
	return nil
}

func validAssetStatus(status string) bool {
	switch status {
	case AssetStatusPurchased, AssetStatusRacked, AssetStatusRunning, AssetStatusMaintenance, AssetStatusRetired:
		return true
	default:
		return false
	}
}
