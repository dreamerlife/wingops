CREATE TABLE IF NOT EXISTS model_groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS models (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES model_groups(id),
  name TEXT NOT NULL,
  display_name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(group_id, name)
);

CREATE TABLE IF NOT EXISTS model_fields (
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

CREATE TABLE IF NOT EXISTS model_relations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
  target_model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
  relation_type TEXT NOT NULL,
  display_name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(source_model_id, target_model_id, relation_type)
);

INSERT INTO permissions (code, description) VALUES
  ('cmdb.model.read', '查看 CMDB 模型'),
  ('cmdb.model.write', '管理 CMDB 模型')
ON CONFLICT (code) DO UPDATE SET description = EXCLUDED.description;
