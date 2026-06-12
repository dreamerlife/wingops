package database

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"wingops/internal/auth"
)

func Seed(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("postgres db is required")
	}
	hash, err := auth.HashPassword("admin123")
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := seedRBAC(tx); err != nil {
			return err
		}
		if err := tx.Exec(`
INSERT INTO users (id, username, password_hash, display_name, status)
VALUES ('00000000-0000-0000-0000-000000000001', 'admin', ?, '管理员', 'active')
ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash, display_name = EXCLUDED.display_name, status = EXCLUDED.status, updated_at = now()
`, hash).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.username = 'admin' AND r.name = 'system_admin'
ON CONFLICT DO NOTHING`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
INSERT INTO system_configs (key, value)
VALUES ('platform.name', to_jsonb('WingOps'::text))
ON CONFLICT (key) DO NOTHING`).Error; err != nil {
			return err
		}
		if err := seedCMDB(tx); err != nil {
			return err
		}
		return tx.Exec(`
INSERT INTO api_keys (name, key_id, secret_hash, status)
VALUES ('本地开发同步 Key', 'dev-sync-key', 'dev-sync-secret', 'active')
ON CONFLICT (key_id) DO NOTHING`).Error
	})
}

func seedRBAC(tx *gorm.DB) error {
	permissions := map[string]string{
		"cmdb.asset.read":     "查看 CMDB 资产",
		"cmdb.asset.write":    "管理 CMDB 资产",
		"cmdb.model.read":     "查看 CMDB 模型",
		"cmdb.model.write":    "管理 CMDB 模型",
		"auth.user.read":      "查看用户",
		"auth.user.write":     "管理用户",
		"auth.role.read":      "查看角色",
		"auth.role.write":     "管理角色和授权",
		"audit.log.read":      "查看审计日志",
		"system.config.read":  "查看系统配置",
		"system.config.write": "管理系统配置",
		"cmdb.apikey.read":    "查看 CMDB API Key",
		"cmdb.apikey.write":   "管理 CMDB API Key",
	}
	for code, description := range permissions {
		if err := tx.Exec(`
INSERT INTO permissions (code, description)
VALUES (?, ?)
ON CONFLICT (code) DO UPDATE SET description = EXCLUDED.description`, code, description).Error; err != nil {
			return err
		}
	}
	rolePermissions := map[string][]string{
		"system_admin": {"cmdb.asset.read", "cmdb.asset.write", "cmdb.model.read", "cmdb.model.write", "auth.user.read", "auth.user.write", "auth.role.read", "auth.role.write", "audit.log.read", "system.config.read", "system.config.write", "cmdb.apikey.read", "cmdb.apikey.write"},
		"ops_admin":    {"cmdb.asset.read", "cmdb.asset.write", "cmdb.model.read", "cmdb.model.write", "audit.log.read", "cmdb.apikey.read", "cmdb.apikey.write"},
		"ops_operator": {"cmdb.asset.read", "cmdb.asset.write"},
		"readonly":     {"cmdb.asset.read", "cmdb.model.read"},
	}
	for role, codes := range rolePermissions {
		for _, code := range codes {
			if err := tx.Exec(`
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = ? AND p.code = ?
ON CONFLICT DO NOTHING`, role, code).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedCMDB(tx *gorm.DB) error {
	groups := []struct {
		name        string
		displayName string
		description string
	}{
		{"compute", "计算资源", "物理服务器、虚拟机和业务应用运行节点"},
		{"network", "网络设备", "交换机、路由器和网络链路设备"},
		{"middleware", "中间件", "Nginx、MySQL、Redis 等基础软件"},
		{"application", "业务应用", "业务系统与应用服务"},
	}
	for _, group := range groups {
		if err := tx.Exec(`
INSERT INTO model_groups (name, display_name, description)
VALUES (?, ?, ?)
ON CONFLICT (name) DO UPDATE SET display_name = EXCLUDED.display_name, description = EXCLUDED.description, updated_at = now()`,
			group.name, group.displayName, group.description).Error; err != nil {
			return err
		}
	}

	models := []struct {
		groupName   string
		name        string
		displayName string
		fields      []map[string]any
	}{
		{"compute", "physical_server", "物理服务器", append(defaultAssetFields("序列号"), numberField("cpu_cores", "CPU 核数"), numberField("memory_gb", "内存 GB"), textField("rack", "机柜位置", false))},
		{"compute", "virtual_machine", "虚拟机", append(defaultAssetFields("实例 ID"), numberField("vcpu", "vCPU"), numberField("memory_gb", "内存 GB"), textField("cluster", "所属集群", false))},
		{"network", "switch", "交换机", append(defaultAssetFields("设备 SN"), numberField("port_count", "端口数"), textField("firmware", "固件版本", false))},
		{"network", "router", "路由器", append(defaultAssetFields("设备 SN"), numberField("interface_count", "接口数"), textField("firmware", "固件版本", false))},
		{"middleware", "nginx", "Nginx", append(defaultAssetFields("服务名"), numberField("listen_port", "监听端口"), textField("config_path", "配置路径", false))},
		{"middleware", "mysql", "MySQL", append(defaultAssetFields("实例名"), numberField("port", "端口"), enumField("role", "角色", []string{"primary", "replica"}))},
		{"middleware", "redis", "Redis", append(defaultAssetFields("实例名"), numberField("port", "端口"), enumField("role", "角色", []string{"primary", "replica", "sentinel"}))},
		{"application", "business_app", "业务应用", append(defaultAssetFields("应用编码"), textField("team", "所属团队", false), enumField("sla_level", "SLA 等级", []string{"gold", "silver", "bronze"}))},
	}
	for _, model := range models {
		var modelID string
		if err := tx.Raw(`
WITH upsert AS (
  INSERT INTO models (group_id, name, display_name, description)
  SELECT id, ?, ?, '' FROM model_groups WHERE name = ?
  ON CONFLICT (group_id, name) DO UPDATE SET display_name = EXCLUDED.display_name, updated_at = now()
  RETURNING id
)
SELECT id FROM upsert`, model.name, model.displayName, model.groupName).Scan(&modelID).Error; err != nil {
			return err
		}
		for index, field := range model.fields {
			options, _ := json.Marshal(field["options"])
			if err := tx.Exec(`
INSERT INTO model_fields (model_id, name, display_name, field_type, required, unique_value, options, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?::jsonb, ?)
ON CONFLICT (model_id, name) DO UPDATE SET display_name = EXCLUDED.display_name, field_type = EXCLUDED.field_type, required = EXCLUDED.required, unique_value = EXCLUDED.unique_value, options = EXCLUDED.options, sort_order = EXCLUDED.sort_order`,
				modelID, field["name"], field["display_name"], field["field_type"], field["required"], field["unique_value"], string(options), index+1).Error; err != nil {
				return err
			}
		}
	}
	if err := seedAssetGroups(tx); err != nil {
		return err
	}
	return seedModelRelations(tx)
}

func seedAssetGroups(tx *gorm.DB) error {
	groups := []AssetGroupSeed{
		{"biz-platform", "平台业务线", "business"},
		{"biz-finance", "财务业务线", "business"},
		{"env-prod", "生产环境", "environment"},
		{"env-staging", "预发环境", "environment"},
		{"region-cn-central", "长沙中心机房", "region"},
	}
	for _, group := range groups {
		if err := tx.Exec(`
INSERT INTO asset_groups (name, display_name, dimension)
VALUES (?, ?, ?)
ON CONFLICT (name) DO UPDATE SET display_name = EXCLUDED.display_name, dimension = EXCLUDED.dimension`,
			group.Name, group.DisplayName, group.Dimension).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedModelRelations(tx *gorm.DB) error {
	relations := []struct {
		source string
		target string
		typ    string
		name   string
	}{
		{"physical_server", "mysql", "runs", "运行"},
		{"physical_server", "redis", "runs", "运行"},
		{"physical_server", "nginx", "runs", "运行"},
		{"business_app", "mysql", "depends_on", "依赖"},
		{"business_app", "redis", "depends_on", "依赖"},
		{"business_app", "nginx", "exposes_by", "入口"},
	}
	for _, relation := range relations {
		if err := tx.Exec(`
INSERT INTO model_relations (source_model_id, target_model_id, relation_type, display_name)
SELECT s.id, t.id, ?, ?
FROM models s, models t
WHERE s.name = ? AND t.name = ?
ON CONFLICT DO NOTHING`, relation.typ, relation.name, relation.source, relation.target).Error; err != nil {
			return err
		}
	}
	return nil
}

type AssetGroupSeed struct {
	Name        string
	DisplayName string
	Dimension   string
}

func defaultAssetFields(uniqueLabel string) []map[string]any {
	return []map[string]any{
		{"name": "name", "display_name": "名称", "field_type": "text", "required": true, "unique_value": false, "options": []string{}},
		{"name": "unique_key_label", "display_name": uniqueLabel, "field_type": "text", "required": false, "unique_value": false, "options": []string{}},
		{"name": "management_ip", "display_name": "管理 IP", "field_type": "ip", "required": true, "unique_value": false, "options": []string{}},
		{"name": "environment", "display_name": "环境", "field_type": "enum", "required": true, "unique_value": false, "options": []string{"prod", "staging", "test", "dev"}},
		{"name": "owner", "display_name": "负责人", "field_type": "text", "required": false, "unique_value": false, "options": []string{}},
	}
}

func textField(name string, displayName string, required bool) map[string]any {
	return map[string]any{"name": name, "display_name": displayName, "field_type": "text", "required": required, "unique_value": false, "options": []string{}}
}

func numberField(name string, displayName string) map[string]any {
	return map[string]any{"name": name, "display_name": displayName, "field_type": "number", "required": false, "unique_value": false, "options": []string{}}
}

func enumField(name string, displayName string, options []string) map[string]any {
	return map[string]any{"name": name, "display_name": displayName, "field_type": "enum", "required": false, "unique_value": false, "options": options}
}
