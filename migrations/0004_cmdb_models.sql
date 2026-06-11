CREATE TABLE model_groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE models (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES model_groups(id),
  name TEXT NOT NULL,
  display_name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(group_id, name)
);

CREATE TABLE model_fields (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  model_id UUID NOT NULL REFERENCES models(id),
  name TEXT NOT NULL,
  display_name TEXT NOT NULL,
  field_type TEXT NOT NULL,
  required BOOLEAN NOT NULL DEFAULT false,
  unique_value BOOLEAN NOT NULL DEFAULT false,
  options JSONB NOT NULL DEFAULT '{}'::jsonb,
  sort_order INT NOT NULL DEFAULT 0,
  UNIQUE(model_id, name)
);

INSERT INTO permissions (code, description) VALUES
  ('cmdb.model.read', '查看 CMDB 模型'),
  ('cmdb.model.write', '管理 CMDB 模型');
