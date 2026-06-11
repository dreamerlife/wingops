import { apiClient } from './client'

export interface ModelGroup {
  id: string
  name: string
  display_name: string
  description: string
}

export interface FieldDefinition {
  name: string
  display_name: string
  field_type: string
  required: boolean
  unique_value: boolean
  options?: string[]
  sort_order: number
}

export interface Model {
  id: string
  group_id: string
  name: string
  display_name: string
  description: string
  fields: FieldDefinition[]
}

export async function listModelGroups(): Promise<ModelGroup[]> {
  const response = await apiClient.get<{ code: number; data: ModelGroup[] }>('/v1/cmdb/model-groups')
  return response.data.data
}

export async function getModel(id: string): Promise<Model> {
  const response = await apiClient.get<{ code: number; data: Model }>(`/v1/cmdb/models/${id}`)
  return response.data.data
}
