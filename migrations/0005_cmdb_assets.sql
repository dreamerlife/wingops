CREATE TABLE IF NOT EXISTS asset_groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  dimension TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  model_id UUID NOT NULL REFERENCES models(id),
  unique_key TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'running',
  attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(model_id, unique_key)
);

CREATE TABLE IF NOT EXISTS asset_group_members (
  asset_id UUID NOT NULL REFERENCES assets(id),
  group_id UUID NOT NULL REFERENCES asset_groups(id),
  PRIMARY KEY(asset_id, group_id)
);

CREATE TABLE IF NOT EXISTS asset_change_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  asset_id UUID NOT NULL REFERENCES assets(id),
  actor_id UUID,
  before_value JSONB NOT NULL,
  after_value JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
