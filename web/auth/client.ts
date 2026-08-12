export type FTNService = {
  id: string;
  name: string;
  enabled: boolean;
};

export type FTNSession = {
  identity_id: string;
  services: FTNService[];
};

async function request<T>(input: RequestInfo, init?: RequestInit): Promise<T> {
  const response = await fetch(input, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
    ...init,
  });

  if (!response.ok) {
    throw new Error(`FTN_AUTH_${response.status}`);
  }
  return response.status === 204 ? (undefined as T) : response.json();
}

export function signIn(login: string, password: string) {
  return request<{ identity_id: string }>("/api/v1/auth/sign-in", {
    method: "POST",
    body: JSON.stringify({ login, password }),
  });
}

export function getSession() {
  return request<FTNSession>("/api/v1/auth/session");
}

export function signOut() {
  return request<void>("/api/v1/auth/sign-out", { method: "POST" });
}
