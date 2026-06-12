import { apiClient } from './client'

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginResult {
  access_token: string
  token_type: string
}

export async function login(payload: LoginPayload): Promise<LoginResult> {
  const response = await apiClient.post<{ code: number; data: LoginResult }>('/v1/auth/login', payload)
  return response.data.data
}
