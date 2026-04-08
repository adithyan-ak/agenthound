import { api } from "./client";
import type {
  LoginResponse,
  CreateTokenResponse,
  APIToken,
} from "./types";

export async function login(
  username: string,
  password: string,
): Promise<LoginResponse> {
  return api.post("auth/login", { json: { username, password } }).json();
}

export async function createToken(
  name: string,
  expiresAt?: string,
): Promise<CreateTokenResponse> {
  const body: Record<string, string> = { name };
  if (expiresAt) {
    body["expires_at"] = expiresAt;
  }
  return api.post("auth/tokens", { json: body }).json();
}

export async function listTokens(): Promise<APIToken[]> {
  return api.get("auth/tokens").json();
}

export async function deleteToken(id: string): Promise<void> {
  await api.delete(`auth/tokens/${id}`);
}
