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
  relations: ModelRelation[]
}

export interface ModelRelation {
  id?: string
  source_model_id?: string
  target_model_id: string
  relation_type: string
  display_name: string
}

export interface Asset {
  id: string
  model_id: string
  unique_key: string
  status: string
  group_ids?: string[]
  attributes: Record<string, unknown>
}

export interface AssetGroup {
  id: string
  name: string
  display_name: string
  dimension: string
}

export interface AssetListResult {
  items: Asset[]
  total: number
  page: number
  page_size: number
}

export interface AssetListQuery {
  model_id?: string
  group_id?: string
  status?: string
  keyword?: string
  page?: number
  page_size?: number
}

export interface AssetChangeLog {
  id: string
  asset_id: string
  actor_id: string
  before_value: Record<string, unknown>
  after_value: Record<string, unknown>
  created_at: string
}

export interface ApiKey {
  id: string
  name: string
  key_id: string
  secret?: string
  status: string
}

interface ApiResponse<T> {
  code: number
  data: T
}

export async function listModelGroups(): Promise<ModelGroup[]> {
  const response = await apiClient.get<ApiResponse<ModelGroup[]>>('/v1/cmdb/model-groups')
  return response.data.data
}

export async function createModelGroup(payload: Omit<ModelGroup, 'id'>): Promise<ModelGroup> {
  const response = await apiClient.post<ApiResponse<ModelGroup>>('/v1/cmdb/model-groups', payload)
  return response.data.data
}

export async function listModels(groupId: string): Promise<Model[]> {
  const response = await apiClient.get<ApiResponse<Model[]>>(`/v1/cmdb/model-groups/${groupId}/models`)
  return response.data.data
}

export async function createModel(groupId: string, payload: Omit<Model, 'id' | 'group_id'>): Promise<Model> {
  const response = await apiClient.post<ApiResponse<Model>>(`/v1/cmdb/model-groups/${groupId}/models`, payload)
  return response.data.data
}

export async function getModel(id: string): Promise<Model> {
  const response = await apiClient.get<ApiResponse<Model>>(`/v1/cmdb/models/${id}`)
  return response.data.data
}

export async function updateModel(id: string, payload: Model): Promise<Model> {
  const response = await apiClient.put<ApiResponse<Model>>(`/v1/cmdb/models/${id}`, payload)
  return response.data.data
}

export async function deleteModel(id: string): Promise<void> {
  await apiClient.delete(`/v1/cmdb/models/${id}`)
}

export async function listAssetGroups(): Promise<AssetGroup[]> {
  const response = await apiClient.get<ApiResponse<AssetGroup[]>>('/v1/cmdb/asset-groups')
  return response.data.data
}

export async function createAssetGroup(payload: Omit<AssetGroup, 'id'>): Promise<AssetGroup> {
  const response = await apiClient.post<ApiResponse<AssetGroup>>('/v1/cmdb/asset-groups', payload)
  return response.data.data
}

export async function listAssets(query: AssetListQuery = {}): Promise<AssetListResult> {
  const response = await apiClient.get<ApiResponse<AssetListResult>>('/v1/cmdb/assets', { params: query })
  return response.data.data
}

export async function createAsset(payload: Omit<Asset, 'id'>): Promise<Asset> {
  const response = await apiClient.post<ApiResponse<Asset>>('/v1/cmdb/assets', payload)
  return response.data.data
}

export async function getAsset(id: string): Promise<Asset> {
  const response = await apiClient.get<ApiResponse<Asset>>(`/v1/cmdb/assets/${id}`)
  return response.data.data
}

export async function updateAsset(id: string, payload: Asset): Promise<Asset> {
  const response = await apiClient.put<ApiResponse<Asset>>(`/v1/cmdb/assets/${id}`, payload)
  return response.data.data
}

export async function deleteAsset(id: string): Promise<void> {
  await apiClient.delete(`/v1/cmdb/assets/${id}`)
}

export async function listAssetChangeLogs(id: string): Promise<AssetChangeLog[]> {
  const response = await apiClient.get<ApiResponse<AssetChangeLog[]>>(`/v1/cmdb/assets/${id}/change-logs`)
  return response.data.data
}

export async function previewCSVImport(file: File) {
  const text = await file.text()
  const response = await apiClient.post<ApiResponse<Array<{ unique_key: string; attributes: Record<string, unknown> }>>>(
    '/v1/cmdb/assets/import/preview',
    text,
    { headers: { 'Content-Type': 'text/csv' } }
  )
  return response.data.data
}

export async function listApiKeys(): Promise<ApiKey[]> {
  const response = await apiClient.get<ApiResponse<ApiKey[]>>('/v1/cmdb/api-keys')
  return response.data.data
}

export async function createApiKey(payload: Pick<ApiKey, 'name'>): Promise<ApiKey> {
  const response = await apiClient.post<ApiResponse<ApiKey>>('/v1/cmdb/api-keys', payload)
  return response.data.data
}

export async function revokeApiKey(id: string): Promise<void> {
  await apiClient.delete(`/v1/cmdb/api-keys/${id}`)
}
