CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL
);

CREATE TABLE permissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL
);

CREATE TABLE user_roles (
  user_id UUID NOT NULL REFERENCES users(id),
  role_id UUID NOT NULL REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
  role_id UUID NOT NULL REFERENCES roles(id),
  permission_id UUID NOT NULL REFERENCES permissions(id),
  PRIMARY KEY (role_id, permission_id)
);

INSERT INTO roles (name, display_name) VALUES
  ('system_admin', '系统管理员'),
  ('ops_admin', '运维管理员'),
  ('ops_operator', '运维操作员'),
  ('readonly', '只读用户');

INSERT INTO permissions (code, description) VALUES
  ('cmdb.asset.read', '查看 CMDB 资产'),
  ('cmdb.asset.write', '管理 CMDB 资产'),
  ('auth.user.read', '查看用户'),
  ('auth.role.read', '查看角色');
