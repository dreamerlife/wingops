package cmdb

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeEnum     FieldType = "enum"
	FieldTypeDate     FieldType = "date"
	FieldTypeIP       FieldType = "ip"
	FieldTypeRelation FieldType = "relation"
)

type ModelGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type Model struct {
	ID          string            `json:"id"`
	GroupID     string            `json:"group_id"`
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Description string            `json:"description"`
	Fields      []FieldDefinition `json:"fields"`
	Relations   []ModelRelation   `json:"relations"`
}

type FieldDefinition struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Type        FieldType `json:"field_type"`
	Required    bool      `json:"required"`
	UniqueValue bool      `json:"unique_value"`
	Options     []string  `json:"options,omitempty"`
	SortOrder   int       `json:"sort_order"`
}

type ModelRelation struct {
	ID            string `json:"id"`
	SourceModelID string `json:"source_model_id"`
	TargetModelID string `json:"target_model_id"`
	RelationType  string `json:"relation_type"`
	DisplayName   string `json:"display_name"`
}

func (f FieldDefinition) Validate(value any) error {
	if value == nil || value == "" {
		if f.Required {
			return fmt.Errorf("%s is required", f.Name)
		}
		return nil
	}

	switch f.Type {
	case FieldTypeIP:
		text, ok := value.(string)
		if !ok || net.ParseIP(text) == nil {
			return fmt.Errorf("%s must be a valid ip", f.Name)
		}
	case FieldTypeEnum:
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s must be an enum value", f.Name)
		}
		for _, option := range f.Options {
			if option == text {
				return nil
			}
		}
		return fmt.Errorf("%s must be one of configured options", f.Name)
	case FieldTypeNumber:
		if !isNumber(value) {
			return fmt.Errorf("%s must be a number", f.Name)
		}
	case FieldTypeText, FieldTypeDate, FieldTypeRelation:
		return nil
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func isNumber(value any) bool {
	switch typed := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(typed, 64)
		return err == nil
	default:
		return false
	}
}
