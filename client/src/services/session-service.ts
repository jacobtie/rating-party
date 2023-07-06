import { baseUrl } from './utils';

type SessionResponse = {
  jwt: string
  isAdmin?: boolean
  gameId?: string
}

export async function signin(username: string, passcode: string): Promise<SessionResponse> {
  const response = await fetch(`${baseUrl}/signin`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, passcode }),
  });
  const session: SessionResponse = await response.json();
  return session;
}
