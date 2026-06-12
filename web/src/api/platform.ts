import { apiClient } from './client'

interface ApiResponse<T> {
  code: number
  data: T
}

export interface User {
  id: string
  username: string
  display_name: string
  status: string
  roles: Array<{ name: string; display_name?: string }>
}

export interface Role {
  name: string
  display_name: string
  permissions: Permission[]
}

export interface Permission {
  code: string
  description: string
}

export interface AuditLog {
  actor_id: string
  method: string
  path: string
  status_code: number
  resource: string
  created_at: string
}

export interface SystemConfig {
  key: string
  value: string
}

export async function listUsers(): Promise<User[]> {
  const response = await apiClient.get<ApiResponse<User[]>>('/v1/auth/users')
  return response.data.data
}

export async function listRoles(): Promise<Role[]> {
  const response = await apiClient.get<ApiResponse<Role[]>>('/v1/auth/roles')
  return response.data.data
}

export async function createUser(payload: {
  username: string
  password: string
  display_name: string
  status: string
  role_names: string[]
}): Promise<User> {
  const response = await apiClient.post<ApiResponse<User>>('/v1/auth/users', payload)
  return response.data.data
}

export async function updateUser(
  id: string,
  payload: {
    username: string
    password?: string
    display_name: string
    status: string
    role_names: string[]
  }
): Promise<User> {
  const response = await apiClient.put<ApiResponse<User>>(`/v1/auth/users/${id}`, payload)
  return response.data.data
}

export async function deleteUser(id: string): Promise<void> {
  await apiClient.delete(`/v1/auth/users/${id}`)
}

export async function listPermissions(): Promise<Permission[]> {
  const response = await apiClient.get<ApiResponse<Permission[]>>('/v1/auth/permissions')
  return response.data.data
}

export async function createRole(payload: {
  name: string
  display_name: string
  permission_codes: string[]
}): Promise<Role> {
  const response = await apiClient.post<ApiResponse<Role>>('/v1/auth/roles', payload)
  return response.data.data
}

export async function updateRole(
  name: string,
  payload: {
    name: string
    display_name: string
    permission_codes: string[]
  }
): Promise<Role> {
  const response = await apiClient.put<ApiResponse<Role>>(`/v1/auth/roles/${name}`, payload)
  return response.data.data
}

export async function deleteRole(name: string): Promise<void> {
  await apiClient.delete(`/v1/auth/roles/${name}`)
}

export async function listAuditLogs(): Promise<AuditLog[]> {
  const response = await apiClient.get<ApiResponse<AuditLog[]>>('/v1/audit/logs')
  return response.data.data
}

export async function listSystemConfigs(): Promise<SystemConfig[]> {
  const response = await apiClient.get<ApiResponse<SystemConfig[]>>('/v1/system/configs')
  return response.data.data
}

export async function saveSystemConfig(payload: SystemConfig): Promise<SystemConfig> {
  const response = await apiClient.put<ApiResponse<SystemConfig>>(`/v1/system/configs/${payload.key}`, payload)
  return response.data.data
}
